package repository

import (
	"gitee.com/cristiane/micro-mall-pay/model/mysql"
	"xorm.io/xorm"
)

func CreatePayRecord(tx *xorm.Session, model *mysql.PayRecord) (err error) {
	_, err = tx.Table(mysql.TablePayRecord).Insert(model)
	return
}
