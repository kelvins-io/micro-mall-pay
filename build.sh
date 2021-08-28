echo 当前分支
git branch

echo 拉取依赖
go mod vendor

cp -n ./etc/app-docker.ini.example ./etc/app.ini

echo 开始构建版本
go build -o micro-mall-pay main.go