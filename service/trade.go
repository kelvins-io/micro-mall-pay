package service

import (
	"context"
	"gitee.com/cristiane/micro-mall-pay/pkg/code"
	"gitee.com/cristiane/micro-mall-pay/proto/micro_mall_pay_proto/pay_business"
	"github.com/google/uuid"
)

func GetTradeUUID(ctx context.Context, req *pay_business.GetTradeUUIDRequest) (tradeUUID string, retCode int)  {
	tradeUUID = uuid.New().String()
	retCode = code.Success
	return
}
