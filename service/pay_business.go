package service

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-pay/model/args"
	"gitee.com/cristiane/micro-mall-pay/model/mysql"
	"gitee.com/cristiane/micro-mall-pay/pkg/code"
	"gitee.com/cristiane/micro-mall-pay/pkg/util"
	"gitee.com/cristiane/micro-mall-pay/proto/micro_mall_pay_proto/pay_business"
	"gitee.com/cristiane/micro-mall-pay/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-pay/repository"
	"gitee.com/cristiane/micro-mall-pay/vars"
	"gitee.com/kelvins-io/common/crypt"
	"gitee.com/kelvins-io/common/errcode"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"strings"
	"time"
	"xorm.io/xorm"
)

func TradePay(ctx context.Context, req *pay_business.TradePayRequest) (payId string, retCode int) {
	retCode = code.Success
	// 参数验证
	if len(req.EntryList) == 0 {
		return
	}
	for i := 0; i < len(req.EntryList); i++ {
		if len(req.EntryList[i].OutTradeNo) == 0 {
			retCode = code.TradeUUIDEmpty
			return
		}
	}
	// 支付状态检查
	retCode = tradePayCheckState(ctx, req)
	if retCode != code.Success {
		return
	}
	// 长事务，多次扣减用户账户在一个事务中完成
	tx := kelvins.XORM_DBEngine.NewSession()
	err := tx.Begin()
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "TradePay Begin err: %v", err)
		retCode = code.ErrorServer
		return
	}
	// 检查用户账户余额
	userAccount, retCode := tradePayCheckUserAccount(ctx, tx, req)
	if retCode != code.Success {
		return
	}

	// 依次支付
	payId = util.GetUUID() // 同一批订单支付交易号唯一
	for i := 0; i < len(req.EntryList); i++ {
		retCode = tradePayOne(ctx, payId, req, i, tx, userAccount)
		if retCode != code.Success {
			return
		}
	}
	// 触发支付事件通知
	retCode = tradeEventNotice(ctx, req, payId)
	if retCode != code.Success {
		errRollback := tx.Rollback()
		if errRollback != nil {
			kelvins.ErrLogger.Errorf(ctx, "TradePay Rollback err: %v", errRollback)
		}
		return
	}
	err = tx.Commit()
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "TradePay Commit err: %v", err)
		retCode = code.ErrorServer
		return
	}

	return
}

func tradeEventNotice(ctx context.Context, req *pay_business.TradePayRequest, payId string) int {
	// 触发支付消息
	pushSer := NewPushNoticeService(vars.TradePayQueueServer, PushMsgTag{
		DeliveryTag:    args.TaskNameTradePayNotice,
		DeliveryErrTag: args.TaskNameTradePayNoticeErr,
		RetryCount:     kelvins.QueueAMQPSetting.TaskRetryCount,
		RetryTimeout:   kelvins.QueueAMQPSetting.TaskRetryTimeout,
	})
	businessMsg := args.CommonBusinessMsg{
		Type: args.TradePayEventTypeCreate,
		Tag:  args.GetMsg(args.TradePayEventTypeCreate),
		UUID: util.GetUUID(),
		Msg: json.MarshalToStringNoError(args.TradePayNotice{
			Uid:    req.OpUid,
			Time:   util.ParseTimeOfStr(time.Now().Unix()),
			PayId:  payId,
			TxCode: req.OutTxCode, // 单次交易号（可能关联多个订单）
		}),
	}
	taskUUID, retCode := pushSer.PushMessage(ctx, businessMsg)
	if retCode != code.Success {
		kelvins.ErrLogger.Errorf(ctx, "trade pay businessMsg: %+v  notice send err: ", businessMsg, errcode.GetErrMsg(retCode))
	} else {
		kelvins.BusinessLogger.Infof(ctx, "trade pay businessMsg businessMsg: %+v  taskUUID :%v", businessMsg, taskUUID)
	}
	return retCode
}

func tradePayOne(ctx context.Context, payId string, req *pay_business.TradePayRequest, i int, tx *xorm.Session, userAccount *mysql.Account) (retCode int) {
	retCode = code.Success
	defer func() {
		if retCode != code.Success {
			errRollback := tx.Rollback()
			if errRollback != nil {
				kelvins.ErrLogger.Errorf(ctx, "tradePayOne Rollback err: %v", errRollback)
			}
		}
	}()
	// 生成支付记录
	payRecord := mysql.PayRecord{
		TxId:        payId,
		OutTradeNo:  req.EntryList[i].OutTradeNo,
		NotifyUrl:   req.EntryList[i].NotifyUrl,
		Description: req.EntryList[i].Description,
		Merchant:    req.EntryList[i].Merchant,
		Attach:      req.EntryList[i].Attach,
		User:        req.Account,
		Amount:      req.EntryList[i].Detail.Amount,
		Reduction:   req.EntryList[i].Detail.Reduction,
		CoinType:    int(req.CoinType),
		PayType:     1,
		PayState:    3,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	err := repository.CreatePayRecord(tx, &payRecord)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CreatePayRecord err: %v, payRecord: %v", err, payRecord)
		retCode = code.ErrorServer
		return
	}
	reqAmount, err := decimal.NewFromString(req.EntryList[i].Detail.Amount)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "tradePayOne NewFromString err: %v, amount: %v", err, req.EntryList[i].Detail.Amount)
		retCode = code.DecimalParseErr
		return
	}
	reduction, err := decimal.NewFromString(req.EntryList[i].Detail.Reduction)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "tradePayOne NewFromString err: %v, amount: %v", err, req.EntryList[i].Detail.Reduction)
		retCode = code.DecimalParseErr
		return
	}
	userBalance, err := decimal.NewFromString(userAccount.Balance)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "TradePay NewFromString err: %v, number: %v", err, userAccount.Balance)
		retCode = code.DecimalParseErr
		return
	}
	merchantAccount, err := repository.GetAccountByTx(tx, sqlSelectCheckUserAccount, req.EntryList[i].Merchant, args.AccountTypeCompany, int(req.CoinType))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetAccountByTx err: %v, owner: %v", err, req.EntryList[i].Merchant)
		retCode = code.ErrorServer
		return
	}
	if merchantAccount.Owner == "" {
		retCode = code.MerchantAccountNotExist
		return
	}
	if merchantAccount.State != 3 {
		retCode = code.MerchantAccountStateLock
		return
	}
	merchantBalance, err := decimal.NewFromString(merchantAccount.Balance)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetAccountByTx err: %v, owner: %v", err, merchantAccount.Balance)
		retCode = code.ErrorServer
		return
	}
	// 生成交易流水
	fromBalance := util.DecimalSub(userBalance, util.DecimalSub(reqAmount, reduction))
	toBalance := util.DecimalAdd(merchantBalance, util.DecimalSub(reqAmount, reduction))
	transaction := mysql.Transaction{
		FromAccountCode: userAccount.AccountCode,
		FromBalance:     fromBalance.String(),
		ToAccountCode:   merchantAccount.AccountCode,
		ToBalance:       toBalance.String(),
		Amount:          util.DecimalSub(reqAmount, reduction).String(),
		Meta:            req.EntryList[i].Description,
		Scene:           req.EntryList[i].Description,
		OpUid:           req.OpUid,
		OpIp:            req.OpIp,
		TxId:            payId,
		Fingerprint:     "",
		PayType:         0,
		PayDesc:         "交易支付",
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}
	transaction.Fingerprint = genTransactionFingerprint(&transaction)
	err = repository.CreateTransaction(tx, &transaction)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CreateTransaction err: %v, transaction: %+v", err, transaction)
		retCode = code.ErrorServer
		return
	}
	lastTxId := uuid.New().String()
	// 扣减用余额，增加商余额
	whereUserAccount := map[string]interface{}{
		"owner":      userAccount.Owner,
		"balance":    userAccount.Balance,
		"last_tx_id": userAccount.LastTxId, // 防止更新期间账户变更
		"state":      userAccount.State,
	}
	userAccountChange := map[string]interface{}{
		"balance":     fromBalance.String(),
		"update_time": time.Now(),
		"last_tx_id":  lastTxId, // 记录本次支付事务ID，对标支付记录
	}
	rowsAffected, err := repository.ChangeAccount(tx, whereUserAccount, userAccountChange)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "ChangeAccount err: %v, userAccountQ: %+v, userAccountChange: %+v", err, whereUserAccount, userAccountChange)
		retCode = code.ErrorServer
		return
	}
	// 没有符合条件的数据行，说明没有更新成功
	if rowsAffected != 1 {
		retCode = code.TransactionFailed
		return
	}
	// 更新扣减了余额后的用户账户（userAccount必须可修改）
	userAccount.Balance = fromBalance.String() // 用户账户剩余金额
	userAccount.LastTxId = lastTxId

	// 增加商户账户余额-，增加商户用户余额应该放在事务最后阶段
	whereMerchantAccount := map[string]interface{}{
		"owner":      merchantAccount.Owner,
		"balance":    merchantAccount.Balance,
		"last_tx_id": merchantAccount.LastTxId, // 防止更新期间账户变更
		"state":      merchantAccount.State,
	}
	merchantAccountChange := map[string]interface{}{
		"balance":     toBalance.String(),
		"update_time": time.Now(),
		"last_tx_id":  lastTxId, // 记录本次支付事务ID，对标支付记录
	}
	rowsAffected, err = repository.ChangeAccount(tx, whereMerchantAccount, merchantAccountChange)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "ChangeAccount err: %v, userAccountQ: %+v, userAccountChange: %+v", err, whereMerchantAccount, userAccountChange)
		retCode = code.ErrorServer
		return
	}
	// 没有符合条件的数据行，说明没有更新成功
	if rowsAffected != 1 {
		retCode = code.TransactionFailed
		return
	}

	return
}

// 生成交易指纹
const appKeyTransaction = "ZpONco7fjvGEFgw4ymMX"

func genTransactionFingerprint(transaction *mysql.Transaction) string {
	params := map[string]string{
		"FromAccountCode": transaction.FromAccountCode,
		"FromBalance":     transaction.FromBalance,
		"ToAccountCode":   transaction.ToAccountCode,
		"ToBalance":       transaction.ToBalance,
		"Amount":          transaction.Amount,
		"OpUid":           fmt.Sprintf("%d", transaction.OpUid),
		"OpIp":            transaction.OpIp,
		"TxId":            transaction.TxId,
		"PayType":         fmt.Sprintf("%d", transaction.PayType),
		"CreateTime":      util.ParseTimeOfStr(transaction.CreateTime.UnixNano()),
	}
	return crypt.Md5Sign(params, appKeyTransaction)
}

const sqlSelectCheckUserAccount = "balance,account_code,owner,balance,last_tx_id,state"

func tradePayCheckUserAccount(ctx context.Context, tx *xorm.Session, req *pay_business.TradePayRequest) (userAccount *mysql.Account, retCode int) {
	retCode = code.Success
	defer func() {
		if retCode != code.Success {
			errRollback := tx.Rollback()
			if errRollback != nil {
				kelvins.ErrLogger.Errorf(ctx, "tradePayCheckUserAccount Rollback err: %v", errRollback)
			}
		}
	}()

	var err error
	userAccount, err = repository.GetAccountByTx(tx, sqlSelectCheckUserAccount, req.Account, args.AccountTypePerson, int(req.CoinType))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetAccountByTx err: %v, owner: %v", err, req.Account)
		retCode = code.ErrorServer
		return
	}
	if userAccount.Owner == "" {
		retCode = code.UserAccountNotExist
		return
	}
	if userAccount.State != 3 {
		retCode = code.UserAccountStateLock
		return
	}
	// 检查用户账户余额
	userBalance, err := decimal.NewFromString(userAccount.Balance)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetAccountByTx NewFromString err: %v, number: %v", err, userAccount.Balance)
		retCode = code.DecimalParseErr
		return
	}
	totalAmount, _ := decimal.NewFromString("0")
	for i := 0; i < len(req.EntryList); i++ {
		amount := req.EntryList[i].Detail.Amount
		amountDecimal, err := decimal.NewFromString(amount)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetAccountByTx NewFromString err: %v, amount: %v", err, amount)
			retCode = code.DecimalParseErr
			return
		}
		totalAmount = util.DecimalAdd(totalAmount, amountDecimal)
	}
	if !util.DecimalGreaterThanOrEqual(userBalance, totalAmount) {
		retCode = code.UserAccountNotEnough
		return
	}
	retCode = code.Success
	return
}

func tradePayCheckState(ctx context.Context, req *pay_business.TradePayRequest) (retCode int) {
	retCode = code.Success
	serverName := args.RpcServiceMicroMallUsers
	conn, err := util.GetGrpcClient(serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		retCode = code.ErrorServer
		return
	}
	defer conn.Close()
	serve := users.NewUsersServiceClient(conn)
	r := users.GetUserAccountIdRequest{
		UidList: []int64{req.OpUid},
	}
	rsp, err := serve.GetUserAccountId(ctx, &r)
	if err != nil || rsp.Common.Code != users.RetCode_SUCCESS {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfo %v,err: %v", serverName, err)
		retCode = code.ErrorServer
		return
	}
	if rsp.Common.Code == users.RetCode_USER_NOT_EXIST {
		retCode = code.UserNotExist
		return
	}
	if rsp.InfoList[0].AccountId == "" {
		retCode = code.UserNotExist
		return
	}
	if rsp.InfoList[0].AccountId != req.Account {
		retCode = code.TradeOrderNotMatchUser
		return
	}
	// 参数验证
	outTradeNoList := make([]string, len(req.EntryList))
	for i := 0; i < len(req.EntryList); i++ {
		outTradeNoList[i] = req.EntryList[i].OutTradeNo
	}
	where := map[string]interface{}{
		"user":         rsp.InfoList[0].AccountId,
		"out_trade_no": outTradeNoList,
	}
	payRecordList, _, err := repository.GetPayRecordList("pay_state", where, nil, nil, 0, 0)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetPayRecordList err: %v, outTradeNoList: %v", err, outTradeNoList)
		retCode = code.ErrorServer
		return
	}
	for i := 0; i < len(payRecordList); i++ {
		if payRecordList[i].PayState == 1 {
			retCode = code.TradePayRun
			return
		}
		if payRecordList[i].PayState == 3 {
			retCode = code.TradePaySuccess
			return
		}
	}
	return
}

func CreateAccount(ctx context.Context, req *pay_business.CreateAccountRequest) (accountCode string, retCode int) {
	retCode = code.Success
	accountType := int(req.AccountType) + 1
	exist, err := repository.CheckAccountExist(req.Owner, accountType, int(req.CoinType))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CheckAccountExist err: %v", err)
		retCode = code.ErrorServer
		return
	}
	if exist {
		retCode = code.AccountExist
		return
	}
	accountCode = util.GetUUID()
	tx := kelvins.XORM_DBEngine.NewSession()
	err = tx.Begin()
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CreateAccount err: %v", err)
		retCode = code.ErrorServer
		return
	}
	// 转账记录
	transaction := mysql.Transaction{
		FromAccountCode: "outside",
		FromBalance:     "0",
		ToAccountCode:   req.Owner,
		ToBalance:       req.Balance,
		Amount:          req.Balance,
		Meta:            "初始账户",
		Scene:           "初始账户",
		OpUid:           0,
		OpIp:            "system",
		TxId:            accountCode,
		Fingerprint:     "",
		PayType:         0,
		PayDesc:         "外部充值",
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}
	transaction.Fingerprint = genTransactionFingerprint(&transaction)
	err = repository.CreateTransaction(tx, &transaction)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			kelvins.ErrLogger.Errorf(ctx, "CreateTransaction Rollback err: %v, transaction: %+v", err, transaction)
		}
		kelvins.ErrLogger.Errorf(ctx, "CreateTransaction err: %v, transaction: %+v", err, transaction)
		retCode = code.ErrorServer
		return
	}
	// 创建账户
	account := mysql.Account{
		AccountCode: accountCode,
		Owner:       req.Owner,
		Balance:     req.Balance,
		CoinType:    int(req.CoinType),
		CoinDesc:    "CNY",
		State:       3,
		AccountType: accountType,
		LastTxId:    accountCode, // 初始值等于AccountCode
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	err = repository.CreateAccount(tx, &account)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			kelvins.ErrLogger.Errorf(ctx, "CreateAccount Rollback err: %v, account: %+v", err, account)
		}
		if strings.Contains(err.Error(), errcode.GetErrMsg(code.DBDuplicateEntry)) {
			retCode = code.AccountExist
			return
		}
		kelvins.ErrLogger.Errorf(ctx, "CreateAccount err: %v, account: %+v", err, account)
		retCode = code.ErrorServer
		return
	}
	err = tx.Commit()
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CreateAccount Commit err: %v", err)
		retCode = code.ErrorServer
		return
	}
	return
}
