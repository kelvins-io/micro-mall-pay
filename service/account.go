package service

import (
	"context"
	"gitee.com/cristiane/micro-mall-pay/model/mysql"
	"gitee.com/cristiane/micro-mall-pay/pkg/code"
	"gitee.com/cristiane/micro-mall-pay/pkg/util"
	"gitee.com/cristiane/micro-mall-pay/proto/micro_mall_pay_proto/pay_business"
	"gitee.com/cristiane/micro-mall-pay/repository"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

const (
	sqlSelectAccount   = "balance,last_tx_id,owner,account_code,state"
	sqlSelectPayRecord = "out_trade_no,pay_state"
)

func AccountCharge(ctx context.Context, req *pay_business.AccountChargeRequest) (retCode int) {
	retCode = code.Success
	if len(req.OutTradeNo) == 0 {
		retCode = code.TradeUUIDEmpty
		return
	}
	// 1 根据外部uuid查找
	wherePayRecord := map[string]interface{}{
		"user":         req.Owner,
		"out_trade_no": req.OutTradeNo,
	}
	payRecordList, err := repository.FindPayRecordList(sqlSelectPayRecord, wherePayRecord)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindPayRecordList err: %v, wherePayRecord: %v", err, wherePayRecord)
		retCode = code.ErrorServer
		return
	}
	for i := 0; i < len(payRecordList); i++ {
		if payRecordList[i].PayState == 1 {
			retCode = code.TradePayRun
			return
		}
		if payRecordList[i].PayState == 3 {
			retCode = code.UserChargeRecordExist
			return
		}
	}
	// 验证账户
	accountList, err := repository.FindAccount(sqlSelectAccount, req.Owner, int(req.AccountType+1), int(req.CoinType))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindAccount err: %v, owner: %v", err, req.Owner)
		retCode = code.ErrorServer
		return
	}
	if len(accountList) == 0 || len(accountList) != len(req.Owner) {
		retCode = code.UserAccountNotExist
		return
	}
	for i := 0; i < len(accountList); i++ {
		if accountList[i].State == 1 {
			retCode = code.UserAccountStateLock
			return
		}
		if accountList[i].State == 2 {
			retCode = code.UserAccountStateInvalid
			return
		}
	}
	ownerToAccount := map[string]mysql.Account{}
	for i := 0; i < len(accountList); i++ {
		ownerToAccount[accountList[i].Owner] = accountList[i]
	}
	tx := kelvins.XORM_DBEngine.NewSession()
	err = tx.Begin()
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "AccountCharge Begin err: %v", err)
		retCode = code.ErrorServer
		return
	}
	defer func() {
		if retCode != code.Success {
			err := tx.Rollback()
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "AccountCharge Rollback err: %v", err)
				return
			}
		}
	}()
	for i := 0; i < len(req.Owner); i++ {
		lastTxId := uuid.New().String()
		account, exist := ownerToAccount[req.Owner[i]]
		if !exist {
			retCode = code.UserAccountNotExist
			return
		}
		account.Balance = decimalZeroCovert(account.Balance)
		req.Amount = decimalZeroCovert(req.Amount)
		transaction := mysql.Transaction{
			FromAccountCode: "outside",
			FromBalance:     "0",
			ToAccountCode:   account.Owner,
			ToBalance:       account.Balance,
			Amount:          req.Amount,
			Meta:            "充值",
			Scene:           "充值",
			OpUid:           req.OpMeta.OpUid,
			OpIp:            req.OpMeta.OpIp,
			TxId:            lastTxId,
			Fingerprint:     "",
			PayType:         0,
			PayDesc:         "外部充值",
			CreateTime:      time.Now(),
			UpdateTime:      time.Now(),
		}
		transaction.Fingerprint = genTransactionFingerprint(&transaction)
		err = repository.CreateTransaction(tx, &transaction)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "CreateTransaction err: %v, transaction: %v", err, json.MarshalToStringNoError(transaction))
			retCode = code.ErrorServer
			return
		}
		balanceOld, err := decimal.NewFromString(account.Balance)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "CreateTransaction balanceOld NewFromString err: %v", err)
			retCode = code.ErrorServer
			return
		}
		balanceDiff, err := decimal.NewFromString(req.Amount)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "CreateTransaction balanceDiff NewFromString err: %v", err)
			retCode = code.ErrorServer
			return
		}
		balanceNew := util.DecimalAdd(balanceOld, balanceDiff)
		// 更新账户余额
		updateAccountWhere := map[string]interface{}{
			"account_code": account.AccountCode, // 账户唯一标识
			"balance":      account.Balance,
			"last_tx_id":   account.LastTxId, // 最后一次操作事务ID
			"state":        account.State,
		}
		updateAccountMaps := map[string]interface{}{
			"last_tx_id": lastTxId,
			"balance":    balanceNew.String(),
		}
		rowsAffected, err := repository.ChangeAccount(tx, updateAccountWhere, updateAccountMaps)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "ChangeAccount err: %v, where: %v, maps: %v", err, json.MarshalToStringNoError(updateAccountWhere), json.MarshalToStringNoError(updateAccountMaps))
			return code.ErrorServer
		}
		// 没有符合条件的数据行，说明没有更新成功
		if rowsAffected != 1 {
			return code.TransactionFailed
		}
	}
	err = tx.Commit()
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "AccountCharge Commit err: %v", err)
		retCode = code.TransactionFailed
		return
	}
	return
}

func FindAccount(ctx context.Context, req *pay_business.FindAccountRequest) (result []*pay_business.AccountEntry, retCode int) {
	result = make([]*pay_business.AccountEntry, 0)
	retCode = code.Success
	accountList, err := repository.FindAccount(sqlSelectAccount, req.Owner, int(req.AccountType+1), int(req.CoinType))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindAccount err: %v, owner: %v", err, req.Owner)
		retCode = code.ErrorServer
		return
	}
	result = make([]*pay_business.AccountEntry, len(accountList))
	for i := 0; i < len(accountList); i++ {
		accountType := pay_business.AccountType_Person
		switch accountList[i].AccountType {
		case 1:
			accountType = pay_business.AccountType_Person
		case 2:
			accountType = pay_business.AccountType_Company
		case 3:
			accountType = pay_business.AccountType_System
		}
		coinType := pay_business.CoinType_CNY
		switch accountList[i].CoinType {
		case 0:
			coinType = pay_business.CoinType_CNY
		case 1:
			coinType = pay_business.CoinType_USD
		}
		entry := &pay_business.AccountEntry{
			Owner:       accountList[i].Owner,
			AccountType: accountType,
			CoinType:    coinType,
			Balance:     accountList[i].Balance,
		}
		result[i] = entry
	}

	return
}
