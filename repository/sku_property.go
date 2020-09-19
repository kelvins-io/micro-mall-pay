package repository

import (
	"gitee.com/cristiane/micro-mall-pay/model/mysql"
	"gitee.com/kelvins-io/kelvins"
	"xorm.io/xorm"
)

// sku 商品属性
func CreateSkuProperty(tx *xorm.Session, model *mysql.SkuProperty) (err error) {
	_, err = tx.Table(mysql.TableSkuProperty).Insert(model)
	return
}

func GetSkuPropertyList(skuCodeList []string) ([]mysql.SkuProperty, error) {
	var result = make([]mysql.SkuProperty, 0)
	err := kelvins.XORM_DBEngine.Table(mysql.TableSkuProperty).In("code", skuCodeList).Find(&result)
	return result, err
}
