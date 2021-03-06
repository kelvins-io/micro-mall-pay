package repository

import (
	"gitee.com/cristiane/micro-mall-pay/model/mysql"
	"gitee.com/kelvins-io/kelvins"
	"xorm.io/xorm"
)

func CheckAccountExist(owner string, accountType, coinType int) (bool, error) {
	var model mysql.Account
	var err error
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableAccount).
		Select("id").
		Where("owner = ?", owner).
		Where("account_type = ?", accountType).
		Where("coin_type = ?", coinType).
		Get(&model)
	if err != nil {
		return false, err
	}
	if model.Id <= 0 {
		return false, nil
	}
	return true, nil
}

func FindAccount(sqlSelect string, ownerList []string, accountType, coinType int) ([]mysql.Account, error) {
	var result = make([]mysql.Account, 0)
	var err error
	session := kelvins.XORM_DBEngine.Table(mysql.TableAccount).
		Select(sqlSelect).
		In("owner", ownerList)
	if accountType != 0 {
		session = session.Where("account_type = ?", accountType)
	}
	err = session.Where("coin_type = ?", coinType).Find(&result)
	return result, err
}

func GetAccountByTx(tx *xorm.Session, sqlSelect string, owner string, accountType, coinType int) (*mysql.Account, error) {
	var model mysql.Account
	var err error
	_, err = tx.Table(mysql.TableAccount).
		Select(sqlSelect).
		Where("owner = ?", owner).
		Where("account_type = ?", accountType).
		Where("coin_type = ?", coinType).
		Get(&model)
	return &model, err
}

func ChangeAccount(tx *xorm.Session, query, maps map[string]interface{}) (int64, error) {
	return tx.Table(mysql.TableAccount).Where(query).Update(maps)
}

func CreateAccount(tx *xorm.Session, model *mysql.Account) (err error) {
	_, err = tx.Table(mysql.TableAccount).Insert(model)
	return err
}
