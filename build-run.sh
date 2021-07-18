git branch

echo 拉取依赖
go mod vendor

echo 开始构建
go build -o micro-mall-pay main.go

cp -n ./etc/app.ini.example ./etc/app.ini
mkdir -p logs

echo 开始运行micro-mall-pay
./micro-mall-pay
