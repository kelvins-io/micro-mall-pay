package repository

import (
	"gitee.com/cristiane/micro-mall-pay/model/mysql"
	"xorm.io/xorm"
)

func CreateTransaction(tx *xorm.Session, model *mysql.Transaction) (err error) {
	_, err = tx.Table(mysql.TableTransaction).Insert(model)
	return
}
