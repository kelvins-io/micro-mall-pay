package server

import (
	"context"
	"gitee.com/cristiane/micro-mall-pay/proto/micro_mall_pay_proto/pay_business"
)

type PayServer struct {
}

func NewPayServer() pay_business.PayBusinessServiceServer {
	return new(PayServer)
}

func (p *PayServer) TradePay(ctx context.Context, req *pay_business.TradePayRequest) (*pay_business.TradePayResponse, error) {
	return &pay_business.TradePayResponse{}, nil
}
