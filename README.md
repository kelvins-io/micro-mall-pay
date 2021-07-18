# micro-mall-pay

#### 介绍
微商城-支付系统

#### 软件架构
grpc

#### 框架，库依赖
kelvins框架支持（gRPC，cron，queue，web支持）：https://gitee.com/kelvins-io/kelvins   
g2cache缓存库支持（两级缓存）：https://gitee.com/kelvins-io/g2cache   

#### 安装教程

1.仅构建  sh build.sh   
2 运行  sh build-run.sh   

#### 使用说明
配置参考
```toml
[kelvins-server]
EndPoint = 8080
IsRecordCallResponse = true

[kelvins-logger]
RootPath = "./logs"
Level = "debug"

[kelvins-mysql]
Host = "127.0.0.1:3306"
UserName = "root"
Password = "xxx"
DBName = "micro_mall_pay"
Charset = "utf8mb4"
PoolNum =  10
MaxIdleConns = 5
ConnMaxLifeSecond = 3600
MultiStatements = true
ParseTime = true

[kelvins-redis]
Host = "127.0.0.1:6379"
Password = "xxx"
DB = 1
PoolNum = 10

[kelvins-queue-amqp]
Broker = "amqp://micro-mall:szJ9aePR@localhost:5672/micro-mall"
DefaultQueue = "trade_pay_notice"
ResultBackend = "redis://xxxx@127.0.0.1:6379/10"
ResultsExpireIn = 36000
Exchange = "trade_pay_notice"
ExchangeType = "direct"
BindingKey = "trade_pay_notice"
PrefetchCount = 5
TaskRetryCount = 3
TaskRetryTimeout = 3600

[email-config]
User = "xxxx@qq.com"
Password = "xxx"
Host = "smtp.qq.com"
Port = "465"

```

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request

