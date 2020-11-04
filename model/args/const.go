package args

const (
	RpcServiceMicroMallUsers = "micro-mall-users"
	RpcServiceMicroMallShop  = "micro-mall-shop"
)

const (
	AccountTypePerson  = 1
	AccountTypeCompany = 2
	AccountTypeSystem  = 3
)

const (
	TaskNameTradePayNotice    = "task_trade_pay_notice"
	TaskNameTradePayNoticeErr = "task_trade_pay_notice_err"
)

type CommonBusinessMsg struct {
	Type int    `json:"type"`
	Tag  string `json:"tag"`
	UUID string `json:"uuid"`
	Msg  string `json:"msg"`
}

type TradePayNotice struct {
	Uid    int64  `json:"uid"`
	Time   string `json:"time"`
	TxCode string `json:"tx_code"`
	PayId  string `json:"pay_id"`
}

const (
	Unknown                 = 0
	TradePayEventTypeCreate = 10016
	TradePayEventTypeExpire = 10017
)

var MsgFlags = map[int]string{
	Unknown:                 "未知",
	TradePayEventTypeCreate: "支付订单创建",
	TradePayEventTypeExpire: "支付订单过期",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Unknown]
}
