package startup

import (
	"gitee.com/cristiane/micro-mall-pay/vars"
	"gitee.com/kelvins-io/kelvins"
	"gitee.com/kelvins-io/kelvins/setup"
)

// SetupVars 加载变量
func SetupVars() error {
	var err error
	vars.TradePayQueueServer, err = setup.NewAMQPQueue(kelvins.QueueAMQPSetting, nil)
	return err
}
