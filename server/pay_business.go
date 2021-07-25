package server

import (
	"context"
	"gitee.com/cristiane/micro-mall-pay/pkg/code"
	"gitee.com/cristiane/micro-mall-pay/proto/micro_mall_pay_proto/pay_business"
	"gitee.com/cristiane/micro-mall-pay/service"
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
	}
	txId, retCode := service.TradePay(ctx, req)
	codeRsp := pay_business.RetCode_SUCCESS
	if retCode != code.Success {
		switch retCode {
		case code.TradeOrderNotMatchUser:
			codeRsp = pay_business.RetCode_TRADE_ORDER_NOT_MATCH_USER
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
		case code.TradePayExpire:
			codeRsp = pay_business.RetCode_TRADE_PAY_EXPIRE
		case code.TradePayRun:
			codeRsp = pay_business.RetCode_TRADE_PAY_RUN
		case code.TradePaySuccess:
			codeRsp = pay_business.RetCode_TRADE_PAY_SUCCESS
		case code.TradeUUIDEmpty:
			result.Common.Code = pay_business.RetCode_TRADE_UUID_EMPTY
		case code.ErrorServer:
			codeRsp = pay_business.RetCode_ERROR
		}
	}
	result.Common.Code = codeRsp
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
	if retCode != code.Success {
		switch retCode {
		case code.AccountExist:
			codeRsp = pay_business.RetCode_ACCOUNT_EXIST
		case code.ErrorServer:
			codeRsp = pay_business.RetCode_ERROR
		}
	}
	result.Common.Code = codeRsp
	result.AccountCode = accountCode
	return &result, nil
}

func (p *PayServer) FindAccount(ctx context.Context, req *pay_business.FindAccountRequest) (*pay_business.FindAccountResponse, error) {
	result := &pay_business.FindAccountResponse{EntryList: nil, Common: &pay_business.CommonResponse{
		Code: pay_business.RetCode_SUCCESS,
	}}
	accountList, retCode := service.FindAccount(ctx, req)
	if retCode != code.Success {
		result.Common.Code = pay_business.RetCode_ERROR
		return result, nil
	}
	result.EntryList = accountList
	return result, nil
}

func (p *PayServer) AccountCharge(ctx context.Context, req *pay_business.AccountChargeRequest) (*pay_business.AccountChargeResponse, error) {
	result := &pay_business.AccountChargeResponse{Common: &pay_business.CommonResponse{
		Code: pay_business.RetCode_SUCCESS,
	}}
	retCode := service.AccountCharge(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.TradePayRun:
			result.Common.Code = pay_business.RetCode_TRADE_PAY_RUN
		case code.TransactionFailed:
			result.Common.Code = pay_business.RetCode_TRANSACTION_FAILED
		case code.UserAccountNotExist:
			result.Common.Code = pay_business.RetCode_USER_ACCOUNT_NOT_EXIST
		case code.UserAccountStateLock:
			result.Common.Code = pay_business.RetCode_USER_ACCOUNT_STATE_LOCK
		case code.UserAccountStateInvalid:
			result.Common.Code = pay_business.RetCode_USER_ACCOUNT_STATE_INVALID
		case code.TradeUUIDEmpty:
			result.Common.Code = pay_business.RetCode_TRADE_UUID_EMPTY
		default:
			result.Common.Code = pay_business.RetCode_ERROR
		}
		return result, nil
	}
	return result, nil
}

func (p *PayServer) GetTradeUUID(ctx context.Context, req *pay_business.GetTradeUUIDRequest) (*pay_business.GetTradeUUIDResponse, error) {
	result := &pay_business.GetTradeUUIDResponse{
		Common: &pay_business.CommonResponse{
			Code: pay_business.RetCode_SUCCESS,
			Msg:  "",
		},
	}
	uuid, retCode := service.GetTradeUUID(ctx, req)
	if retCode != code.Success {
		result.Common.Code = pay_business.RetCode_ERROR
	}
	result.Uuid = uuid
	return result, nil
}
