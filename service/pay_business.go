package service

import (
	"context"
	"gitee.com/cristiane/micro-mall-pay/model/args"
	"gitee.com/cristiane/micro-mall-pay/model/mysql"
	"gitee.com/cristiane/micro-mall-pay/pkg/code"
	"gitee.com/cristiane/micro-mall-pay/pkg/util"
	"gitee.com/cristiane/micro-mall-pay/proto/micro_mall_pay_proto/pay_business"
	"gitee.com/cristiane/micro-mall-pay/repository"
	"gitee.com/cristiane/micro-mall-pay/vars"
	"gitee.com/kelvins-io/common/errcode"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

func TradePay(ctx context.Context, req *pay_business.TradePayRequest) (txId string, retCode int) {
	retCode = code.Success
	// 参数验证
	// 长事务，多次扣减用户账户在一个事务中完成
	tx := kelvins.XORM_DBEngine.NewSession()
	userAccount, err := repository.GetAccountByTx(tx, req.Account, args.AccountTypePerson, int(req.CoinType))
	if err != nil {
		tx.Rollback()
		kelvins.ErrLogger.Errorf(ctx, "GetAccount err: %v, owner: %v", err, req.Account)
		retCode = code.ErrorServer
		return
	}
	if userAccount.Owner == "" {
		tx.Rollback()
		retCode = code.UserAccountNotExist
		return
	}
	if userAccount.State != 3 {
		tx.Rollback()
		retCode = code.UserAccountStateLock
		return
	}
	// 检查用户账户余额
	userBalance, err := decimal.NewFromString(userAccount.Balance)
	if err != nil {
		tx.Rollback()
		kelvins.ErrLogger.Errorf(ctx, "NewFromString err: %v, number: %v", err, userAccount.Balance)
		retCode = code.DecimalParseErr
		return
	}
	totalAmount, _ := decimal.NewFromString("0")
	for i := 0; i < len(req.EntryList); i++ {
		amount := req.EntryList[i].Detail.Amount
		amountDecimal, err := decimal.NewFromString(amount)
		if err != nil {
			tx.Rollback()
			kelvins.ErrLogger.Errorf(ctx, "NewFromString err: %v, amount: %v", err, amount)
			retCode = code.DecimalParseErr
			return
		}
		totalAmount = util.DecimalAdd(totalAmount, amountDecimal)
	}
	if !util.DecimalGreaterThanOrEqual(userBalance, totalAmount) {
		tx.Rollback()
		retCode = code.UserAccountNotEnough
		return
	}

	// 同一批订单支付交易号唯一
	txId = util.GetUUID()
	now := time.Now()
	for i := 0; i < len(req.EntryList); i++ {
		// 生成支付记录
		payRecord := mysql.PayRecord{
			TxId:        txId,
			OutTradeNo:  req.EntryList[i].OutTradeNo,
			TimeExpire:  time.Now().Add(30 * time.Minute),
			NotifyUrl:   req.EntryList[i].NotifyUrl,
			Description: req.EntryList[i].Description,
			Merchant:    req.EntryList[i].Merchant,
			Attach:      req.EntryList[i].Attach,
			User:        req.Account,
			Amount:      req.EntryList[i].Detail.Amount,
			Reduction:   req.EntryList[i].Detail.Reduction,
			CoinType:    int(req.CoinType),
			PayType:     1,
			PayState:    1,
			CreateTime:  now,
			UpdateTime:  now,
		}
		err = repository.CreatePayRecord(tx, &payRecord)
		if err != nil {
			_ = tx.Rollback()
			kelvins.ErrLogger.Errorf(ctx, "CreatePayRecord err: %v, payRecord: %v", err, payRecord)
			retCode = code.ErrorServer
			return
		}
		reqAmount, err := decimal.NewFromString(req.EntryList[i].Detail.Amount)
		if err != nil {
			_ = tx.Rollback()
			kelvins.ErrLogger.Errorf(ctx, "NewFromString err: %v, amount: %v", err, req.EntryList[i].Detail.Amount)
			retCode = code.DecimalParseErr
			return
		}
		reduction, err := decimal.NewFromString(req.EntryList[i].Detail.Reduction)
		if err != nil {
			_ = tx.Rollback()
			kelvins.ErrLogger.Errorf(ctx, "NewFromString err: %v, amount: %v", err, req.EntryList[i].Detail.Reduction)
			retCode = code.DecimalParseErr
			return
		}
		merchantAccount, err := repository.GetAccountByTx(tx, req.EntryList[i].Merchant, args.AccountTypeCompany, int(req.CoinType))
		if err != nil {
			_ = tx.Rollback()
			kelvins.ErrLogger.Errorf(ctx, "GetAccount err: %v, owner: %v", err, req.EntryList[i].Merchant)
			retCode = code.ErrorServer
			return
		}
		if merchantAccount.Owner == "" {
			_ = tx.Rollback()
			retCode = code.MerchantAccountNotExist
			return
		}
		if merchantAccount.State != 3 {
			_ = tx.Rollback()
			retCode = code.MerchantAccountStateLock
			return
		}
		merchantBalance, err := decimal.NewFromString(merchantAccount.Balance)
		if err != nil {
			_ = tx.Rollback()
			kelvins.ErrLogger.Errorf(ctx, "GetAccount err: %v, owner: %v", err, merchantAccount.Balance)
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
			TxId:            txId,
			Fingerprint:     time.Now().String(),
			PayType:         0,
			PayDesc:         "交易支付",
			CreateTime:      now,
			UpdateTime:      now,
		}
		err = repository.CreateTransaction(tx, &transaction)
		if err != nil {
			_ = tx.Rollback()
			kelvins.ErrLogger.Errorf(ctx, "CreateTransaction err: %v, transaction: %+v", err, transaction)
			retCode = code.ErrorServer
			return
		}
		// 扣减用余额，增加商余额
		whereUserAccount := map[string]interface{}{
			"owner":   userAccount.Owner,
			"balance": userAccount.Balance,
		}
		userAccountChange := map[string]interface{}{
			"balance":     fromBalance.String(),
			"update_time": now,
		}
		r, err := repository.ChangeAccount(tx, whereUserAccount, userAccountChange)
		if err != nil {
			_ = tx.Rollback()
			kelvins.ErrLogger.Errorf(ctx, "ChangeAccount err: %v, userAccountQ: %+v, userAccountChange: %+v", err, whereUserAccount, userAccountChange)
			retCode = code.ErrorServer
			return
		}
		// 没有符合条件的数据行，说明没有更新成功
		if r <= 0 {
			_ = tx.Rollback()
			retCode = code.TransactionFailed
			return
		}
		// 更新扣减了余额后的用户账户
		userBalance = fromBalance
		userAccount.Balance = userBalance.String()

		// 增加商户账户余额-，增加商户用户余额应该放在事务最后阶段
		whereMerchantAccount := map[string]interface{}{
			"owner":   merchantAccount.Owner,
			"balance": merchantAccount.Balance,
		}
		merchantAccountChange := map[string]interface{}{
			"balance":     toBalance.String(),
			"update_time": now,
		}
		r, err = repository.ChangeAccount(tx, whereMerchantAccount, merchantAccountChange)
		if err != nil {
			_ = tx.Rollback()
			kelvins.ErrLogger.Errorf(ctx, "ChangeAccount err: %v, userAccountQ: %+v, userAccountChange: %+v", err, whereMerchantAccount, userAccountChange)
			retCode = code.ErrorServer
			return
		}
		// 没有符合条件的数据行，说明没有更新成功
		if r <= 0 {
			_ = tx.Rollback()
			retCode = code.TransactionFailed
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		retCode = code.ErrorServer
		return
	}
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
			TxCode: txId,
		}),
	}
	taskUUID, retCode := pushSer.PushMessage(ctx, businessMsg)
	if retCode != code.Success {
		kelvins.ErrLogger.Errorf(ctx, "trade pay businessMsg: %+v  notice send err: ", businessMsg, errcode.GetErrMsg(retCode))
	} else {
		kelvins.BusinessLogger.Infof(ctx, "trade pay businessMsg businessMsg: %+v  taskUUID :%v", businessMsg, taskUUID)
	}

	return
}

func CreateAccount(ctx context.Context, req *pay_business.CreateAccountRequest) (accountCode string, retCode int) {
	retCode = code.Success
	accountCode = util.GetUUID()
	account := mysql.Account{
		AccountCode: accountCode,
		Owner:       req.Owner,
		Balance:     req.Balance,
		CoinType:    int(req.CoinType),
		CoinDesc:    "CNY",
		State:       3,
		AccountType: int(req.AccountType) + 1,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	err := repository.CreateAccount(&account)
	if err != nil {
		if strings.Contains(err.Error(), errcode.GetErrMsg(code.DBDuplicateEntry)) {
			retCode = code.AccountExist
			return
		}
		kelvins.ErrLogger.Errorf(ctx, "CreateAccount err: %v, account: %+v", err, account)
		retCode = code.ErrorServer
	}
	return
}
