package repository

import (
	"gitee.com/cristiane/micro-mall-pay/model/mysql"
	"gitee.com/kelvins-io/kelvins"
	"xorm.io/xorm"
)

// sku库存
func CreateSkuInventory(tx *xorm.Session, model *mysql.SkuInventory) (err error) {
	_, err = tx.Table(mysql.TableSkuInventory).Insert(model)
	return
}

func CheckSkuInventoryExist(skuCode string) (exist bool, err error) {
	var model mysql.SkuInventory
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableSkuInventory).
		Select("id").
		Where("sku_code = ?", skuCode).Get(&model)
	if err != nil {
		return false, err
	}
	if model.Id > 0 {
		return true, nil
	}
	return false, nil
}

func GetSkuInventoryListByShopId(shopId int64) ([]mysql.SkuInventory, error) {
	var result = make([]mysql.SkuInventory, 0)
	session := kelvins.XORM_DBEngine.Table(mysql.TableSkuInventory)
	if shopId > 0 {
		session = session.Where("shop_id = ?", shopId)
	}
	err := session.Desc("create_time").Find(&result)
	return result, err
}
