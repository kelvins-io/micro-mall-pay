echo 拉取依赖
go mod vendor

cp -f ./etc/app-docker.ini.example ./etc/app.ini

echo 开始构建版本
go build -o micro-mall-pay main.go