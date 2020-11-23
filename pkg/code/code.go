package code

import "gitee.com/kelvins-io/common/errcode"

const (
	Success                  = 29000000
	ErrorServer              = 29000001
	DecimalParseErr          = 29000002
	UserNotExist             = 29000005
	UserExist                = 29000006
	DBDuplicateEntry         = 29000007
	MerchantExist            = 29000008
	MerchantNotExist         = 29000009
	ShopBusinessExist        = 29000010
	ShopBusinessNotExist     = 29000011
	SkuCodeEmpty             = 29000012
	SkuCodeNotExist          = 29000013
	SkuCodeExist             = 29000014
	UserAccountNotExist      = 29000015
	MerchantAccountNotExist  = 29000016
	UserAccountStateLock     = 29000017
	MerchantAccountStateLock = 29000018
	UserAccountNotEnough     = 29000019
	TransactionFailed        = 29000020
	AccountExist             = 29000021
	TradePayRun              = 29000022
	TradePaySuccess          = 29000023
	TradePayExpire           = 29000024
	TradeOrderNotMatchUser   = 29000025
	UserAccountStateInvalid  = 29000026
)

var ErrMap = make(map[int]string)

func init() {
	dict := map[int]string{
		Success:                  "OK",
		ErrorServer:              "服务器错误",
		UserNotExist:             "用户不存在",
		DBDuplicateEntry:         "Duplicate entry",
		DecimalParseErr:          "浮点数解析错误",
		UserExist:                "已存在用户记录，请勿重复创建",
		MerchantExist:            "商户认证材料已存在",
		MerchantNotExist:         "商户未提交材料",
		ShopBusinessExist:        "店铺申请材料已存在",
		ShopBusinessNotExist:     "商户未提交店铺材料",
		SkuCodeEmpty:             "商品唯一code为空",
		SkuCodeNotExist:          "商品唯一code在系统找不到",
		SkuCodeExist:             "商品唯一code已存在系统",
		UserAccountNotExist:      "用户账户不存在",
		MerchantAccountNotExist:  "商户账户不存在",
		UserAccountStateLock:     "用户账户被锁定",
		MerchantAccountStateLock: "商户账户被锁定",
		UserAccountNotEnough:     "用户账户余额不足",
		TransactionFailed:        "交易不成功",
		AccountExist:             "账户已存在",
		TradePayRun:              "订单已在支付中",
		TradePaySuccess:          "订单已经支付",
		TradePayExpire:           "订单支付时间过期",
		TradeOrderNotMatchUser:   "交易订单不匹配当前用户",
		UserAccountStateInvalid:  "用户账户无效",
	}
	errcode.RegisterErrMsgDict(dict)
	for key, _ := range dict {
		ErrMap[key] = dict[key]
	}
}
