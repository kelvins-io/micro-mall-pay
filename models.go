package models

import (
	"time"
)

type Account struct {
	Id          int64     `xorm:"pk autoincr comment('自增ID') BIGINT"`
	AccountCode string    `xorm:"not null comment('账户主键') CHAR(50)"`
	Owner       string    `xorm:"not null comment('账户所有者') unique(account_index) CHAR(36)"`
	Balance     string    `xorm:"comment('账户余额') DECIMAL(64,4)"`
	CoinType    int       `xorm:"not null default 0 comment('币种类型，0-rmb，1-usdt') unique(account_index) TINYINT"`
	CoinDesc    string    `xorm:"comment('币种描述') VARCHAR(64)"`
	State       int       `xorm:"comment('状态，1无效，2锁定，3正常') TINYINT"`
	AccountType int       `xorm:"not null comment('账户类型，1-个人账户，2-公司账户，3-系统账户') unique(account_index) TINYINT"`
	CreateTime  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') index DATETIME"`
	UpdateTime  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
	LastTxId    string    `xorm:"not null default '99' comment('最后一次事务ID') CHAR(40)"`
}

type PayRecord struct {
	Id          int64     `xorm:"pk autoincr comment('自增ID') BIGINT"`
	TxId        string    `xorm:"not null comment('批次交易号') index CHAR(40)"`
	OutTradeNo  string    `xorm:"not null comment('外部商户订单号') index CHAR(40)"`
	NotifyUrl   string    `xorm:"comment('交易结果通知地址') VARCHAR(255)"`
	Description string    `xorm:"comment('交易描述') VARCHAR(255)"`
	Merchant    string    `xorm:"not null comment('交易商户ID') index CHAR(40)"`
	Attach      string    `xorm:"comment('交易留言') TEXT"`
	User        string    `xorm:"not null comment('交易用户ID') index CHAR(40)"`
	Amount      string    `xorm:"not null comment('交易数量') DECIMAL(64,4)"`
	CoinType    int       `xorm:"not null default 0 comment('交易币种，0-cny,1-usd') TINYINT"`
	Reduction   string    `xorm:"comment('满减优惠') DECIMAL(64,4)"`
	PayType     int       `xorm:"not null comment('交易类型，1入账，2退款') TINYINT"`
	PayState    int       `xorm:"comment('支付状态，0-未支付，1-支付中，2-支付失败，3-支付成功') TINYINT"`
	CreateTime  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('修改时间') DATETIME"`
}

type Transaction struct {
	Id              int64     `xorm:"pk autoincr comment('交易ID') BIGINT"`
	FromAccountCode string    `xorm:"not null default '0' comment('转出账户ID') CHAR(36)"`
	FromBalance     string    `xorm:"default 0.0000 comment('转出后账户余额') DECIMAL(64,4)"`
	ToAccountCode   string    `xorm:"not null default '0' comment('转入账户ID') CHAR(36)"`
	ToBalance       string    `xorm:"comment('转入后账户余额') DECIMAL(64,4)"`
	Amount          string    `xorm:"comment('交易金额') DECIMAL(64,4)"`
	Meta            string    `xorm:"comment('转账说明') VARCHAR(255)"`
	Scene           string    `xorm:"comment('支付场景') VARCHAR(64)"`
	OpUid           int64     `xorm:"not null comment('操作用户UID') BIGINT"`
	OpIp            string    `xorm:"comment('操作的IP') VARCHAR(16)"`
	TxId            string    `xorm:"comment('对应交易号') CHAR(36)"`
	Fingerprint     string    `xorm:"not null comment('防篡改指纹') TEXT"`
	PayType         int       `xorm:"default 0 comment('支付方式，0系统操作，1-银行卡，2-信用卡,3-支付宝,4-微信支付,5-京东支付') TINYINT"`
	PayDesc         string    `xorm:"comment('支付方式描述') VARCHAR(36)"`
	CreateTime      time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime      time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('修改时间') DATETIME"`
}
