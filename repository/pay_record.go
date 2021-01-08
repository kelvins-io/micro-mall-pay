package repository

import (
	"gitee.com/cristiane/micro-mall-pay/model/mysql"
	"gitee.com/kelvins-io/kelvins"
	"xorm.io/xorm"
)

func CreatePayRecord(tx *xorm.Session, model *mysql.PayRecord) (err error) {
	_, err = tx.Table(mysql.TablePayRecord).Insert(model)
	return
}

func GetPayRecordList(selectSql string, where interface{}, outTradeNo, orderByAsc, orderByDesc []string, pageSize, pageNum int) ([]mysql.PayRecord, int64, error) {
	var result = make([]mysql.PayRecord, 0)
	session := kelvins.XORM_DBEngine.Table(mysql.TablePayRecord).Select(selectSql).
		In("out_trade_no", outTradeNo).
		Where(where).
		Asc(orderByAsc...).
		Desc(orderByDesc...)
	if pageSize > 0 && pageNum >= 1 {
		session = session.Limit(pageSize, (pageNum-1)*pageSize)
	}
	total, err := session.FindAndCount(&result)
	return result, total, err
}

func FindPayRecordList(selectSql string, where interface{}) ([]mysql.PayRecord,error)  {
	var result = make([]mysql.PayRecord, 0)
	err :=  kelvins.XORM_DBEngine.Table(mysql.TablePayRecord).Select(selectSql).Where(where).Find(&result)
	return result,err
}
