package server

import (
	"context"
	"gitee.com/cristiane/micro-mall-pay/pkg/code"
	"gitee.com/cristiane/micro-mall-pay/proto/micro_mall_pay_proto/pay_business"
	"gitee.com/cristiane/micro-mall-pay/service"
	"gitee.com/kelvins-io/common/errcode"
)

type PayServer struct {
}

func NewPayServer() pay_business.PayBusinessServiceServer {
	return new(PayServer)
}

func (p *PayServer) TradePay(ctx context.Context, req *pay_business.TradePayRequest) (*pay_business.TradePayResponse, error) {
	var result pay_business.TradePayResponse
	result.Common = &pay_business.CommonResponse{
		Code: pay_business.RetCode_SUCCESS,
		Msg:  "ok",
	}

	txId, retCode := service.TradePay(ctx, req)
	codeRsp := pay_business.RetCode_SUCCESS
	msgRsp := errcode.GetErrMsg(retCode)
	if retCode != code.Success {
		switch retCode {
		case code.UserAccountNotExist:
			codeRsp = pay_business.RetCode_USER_ACCOUNT_NOT_EXIST
		case code.UserAccountNotEnough:
			codeRsp = pay_business.RetCode_USER_BALANCE_NOT_ENOUGH
		case code.UserAccountStateLock:
			codeRsp = pay_business.RetCode_USER_ACCOUNT_STATE_LOCK
		case code.MerchantAccountNotExist:
			codeRsp = pay_business.RetCode_MERCHANT_ACCOUNT_NOT_EXIST
		case code.MerchantAccountStateLock:
			codeRsp = pay_business.RetCode_MERCHANT_ACCOUNT_STATE_LOCK
		case code.DecimalParseErr:
			codeRsp = pay_business.RetCode_DECIMAL_PARSE_ERR
		case code.TransactionFailed:
			codeRsp = pay_business.RetCode_TRANSACTION_FAILED
		case code.TradePayRun:
			codeRsp = pay_business.RetCode_TRADE_PAY_RUN
		case code.TradePaySuccess:
			codeRsp = pay_business.RetCode_TRADE_PAY_SUCCESS
		case code.ErrorServer:
			codeRsp = pay_business.RetCode_ERROR
		}
	}
	result.Common.Code = codeRsp
	result.Common.Msg = msgRsp
	result.TradeId = txId

	return &result, nil
}

func (p *PayServer) CreateAccount(ctx context.Context, req *pay_business.CreateAccountRequest) (*pay_business.CreateAccountResponse, error) {
	var result pay_business.CreateAccountResponse
	result.Common = &pay_business.CommonResponse{
		Code: pay_business.RetCode_SUCCESS,
		Msg:  "ok",
	}
	accountCode, retCode := service.CreateAccount(ctx, req)
	codeRsp := pay_business.RetCode_SUCCESS
	msgRsp := errcode.GetErrMsg(retCode)
	if retCode != code.Success {
		switch retCode {
		case code.AccountExist:
			codeRsp = pay_business.RetCode_ACCOUNT_EXIST
		case code.ErrorServer:
			codeRsp = pay_business.RetCode_ERROR
		}
	}
	result.Common.Code = codeRsp
	result.Common.Msg = msgRsp
	result.AccountCode = accountCode
	return &result, nil
}

func (p *PayServer) GetAccount(ctx context.Context, req *pay_business.GetAccountRequest) (*pay_business.GetAccountResponse, error) {
	return &pay_business.GetAccountResponse{}, nil
}
