package args

type SkuInventoryInfo struct {
	SkuCode       string `json:"sku_code"`
	Name          string `json:"name"`
	Price         string `json:"price"`
	Title         string `json:"title"`
	SubTitle      string `json:"sub_title"`
	Desc          string `json:"desc"`
	Production    string `json:"production"`
	Supplier      string `json:"supplier"`
	Category      int32  `json:"category"`
	Color         string `json:"color"`
	ColorCode     int32  `json:"color_code"`
	Specification string `json:"specification"`
	DescLink      string `json:"desc_link"`
	State         int32  `json:"state"`
	// 其他属性...

	Amount int64 `json:"amount"`
	ShopId int64 `json:"shop_id"`
}

type SkuPropertyEx struct {
	OpUid             int64  `bson:"op_uid"`
	OpIp              string `bson:"op_ip"`
	ShopId            int64  `bson:"shop_id"`
	SkuCode           string `bson:"sku_code"`
	Name              string `bson:"name"`
	Size              string `bson:"size"`
	Shape             string `bson:"shape"`
	ProductionCountry string `bson:"production_country"`
	ProductionDate    string `bson:"production_date"`
	ShelfLife         string `bson:"shelf_life"`
}

const (
	RpcServiceMicroMallUsers = "micro-mall-users"
	RpcServiceMicroMallShop  = "micro-mall-shop"
)
