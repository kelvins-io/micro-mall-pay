# micro-mall-pay

#### 介绍
微商城-支付系统

#### 软件架构
grpc

#### 框架，库依赖
kelvins框架支持（gRPC，cron，queue，web支持）：https://gitee.com/kelvins-io/kelvins   
g2cache缓存库支持（两级缓存）：https://gitee.com/kelvins-io/g2cache   

#### 项目问题交流
QQ群： （micro-mall-pay交流群）  
![avatar](./micro-mall-pay.JPG)

#### 如何赞助
开源项目的发展离不开大家的鼓励和赞赏，扫描下方二维码鼓励一下吧   
![avatar](./微信赞赏码.JPG)

支付宝   
![avatar](./支付宝赞赏码.JPG)

#### 赞助商列表
昵称 | 赞助金额 |  时间 | 留言
---|------|------|---
雨化田 | 100元 | 2021-1-25 | 一起加入
thomas | 100元 | 2021-2-18 | 指导
皮卡猪 | 250元 | 2021-2-20 | 支持大佬
*抹 | 20 | 2021-3-19 | 资金有限，支持下
*康 | 66.66 | 2021-4-15 | 加油
Bleem | -goland正版license | 2021-4-18 | 落地验证码限制以及缓存实施
Christible | 66.00 | 2021-4-26 | 大神，膜拜。资金有限
剑峰 | 50.00 | 2021-5-10 | 支持下
mu | 100.00 | 2021-6-9 | 意思意思
osc | -200.00 | 2021-7-9 | 落地docker构建方案
这个杀手有点冷 | 150.00 | 2021-7-11 | 很好的一个项目
pick刘 | 50.00 | 2021-8-3 | 有料
IT詹天佑 | 88.00 | 2021-8-15 | 喜欢这一点代码
*浩 | 20.00 | 2021-8-25 | 请喝一杯奶茶
Doyle | 100.00 | 2021-8-31 | 一点小意思
星辰大海 | 200.00 | 2021-9-5 | 很好的项目
黔驴技穷 | 100.00 | 2021-10-2 | 问题咨询
_天行健_ | 100.00 | 2021-10-15 | 部署咨询
东正 | 20.00 | 2021-10-20 | 咨询解答
Jackson | 20.00 | 2021-10-28 | 部署咨询
農民 | 20.00 | 2021-11-2 | 一杯茶颜悦色
ps | 50.00 | 2021-11-6 | 项目咨询
Mark | 66.00 | 2021-11-14 | etcd搭建协助
Micky | 20.00 | 2021-11-18 | 赞助
井 | 50.00 | 2021-11-19 | 赞助
农民GG | 500.00 | 2021-11-24 | 支持开源项目
Z*k | 20.00 | 2021-11-27 | 支持开源项目
*左 | 20.00 | 2021-11-29 | 赞助
奈何桥 | 300.00 | 2021-12-2 | 运行部署
p神 | 1000.00 | 2021-12-5 | 个人赞助
曹大 | 1000.00 | 2021-12-11 | 个人赞助
*铭 | 50.00 | 2021-12-19 | 赞助
番茄可乐 | 100.00 | 2022-3-4 | 赞助
东窗事发 | 80.00 | 2022-3-10 | 运行支持

#### 安装教程

1.仅构建  sh build.sh   
2 运行  sh build-run.sh   
3 停止 sh stop.sh

#### 使用说明
配置参考
```toml
[kelvins-server]
Environment = "dev"

[kelvins-logger]
RootPath = "./logs"
Level = "debug"

[kelvins-auth]
Token = "c9VW6ForlmzdeDkZE2i8"
TransportSecurity = false
ExpireSecond = 100

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
Enable = false
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

