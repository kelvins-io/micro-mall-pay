package repository

import (
	"context"
	"gitee.com/cristiane/micro-mall-pay/model/args"
	"gitee.com/cristiane/micro-mall-pay/model/mongodb"
	"gitee.com/cristiane/micro-mall-pay/vars"
	"sync"
)

var one sync.Once

func createIndexes() {
	one.Do(func() {
		var uniques = []string{"sku_code,shop_id"}
		var indexes = []string{"shape", "shop_id", "sku_code"}
		vars.MongoDBDatabase.Collection(mongodb.TableSkuPropertyEx).EnsureIndexes(context.Background(), uniques, indexes)
	})
}

func CreateSkuPropertyEx(ctx context.Context, req *args.SkuPropertyEx) (err error) {
	// 创建索引
	createIndexes()
	// 插入记录
	_, err = vars.MongoDBDatabase.Collection(mongodb.TableSkuPropertyEx).InsertOne(ctx, req)
	return
}

func GetSkuPropertyExList(ctx context.Context, query map[string]interface{}) ([]args.SkuPropertyEx, error) {
	var skuExList []args.SkuPropertyEx
	err := vars.MongoDBDatabase.Collection(mongodb.TableSkuPropertyEx).Find(ctx, query).Sort("sku_code").Limit(100).All(&skuExList)
	return skuExList, err
}
