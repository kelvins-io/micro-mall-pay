package repository

import (
	"gitee.com/cristiane/micro-mall-pay/model/mysql"
	"gitee.com/kelvins-io/kelvins"
	"xorm.io/xorm"
)

func GetAccount(owner string, accountType, coinType int) (*mysql.Account, error) {
	var model mysql.Account
	var err error
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableAccount).
		Where("owner = ?", owner).
		Where("account_type = ?", accountType).
		Where("coin_type = ?", coinType).
		Get(&model)
	return &model, err
}

func GetAccountByTx(tx *xorm.Session, owner string, accountType, coinType int) (*mysql.Account, error) {
	var model mysql.Account
	var err error
	_, err = tx.Table(mysql.TableAccount).
		Where("owner = ?", owner).
		Where("account_type = ?", accountType).
		Where("coin_type = ?", coinType).
		Get(&model)
	return &model, err
}

func GetAccountList(ownerList []string, coinType int) ([]mysql.Account, error) {
	var result = make([]mysql.Account, 0)
	err := kelvins.XORM_DBEngine.Table(mysql.TableAccount).
		Where("coin_type = ?", coinType).
		In("owner", ownerList).
		Find(&result)
	return result, err
}

func ChangeAccount(tx *xorm.Session, query, maps map[string]interface{}) (int64, error) {
	return tx.Table(mysql.TableAccount).Where(query).Update(maps)
}

func CreateAccount(model *mysql.Account) (err error) {
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableAccount).Insert(model)
	return err
}
