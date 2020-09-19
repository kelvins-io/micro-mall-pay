package startup

import (
	"gitee.com/cristiane/micro-mall-pay/vars"
	"gitee.com/kelvins-io/kelvins/config"
	"log"
)

const (
	SectionEmailConfig = "email-config"
)

// LoadConfig 加载配置对象映射
func LoadConfig() error {
	// 加载email数据源
	log.Printf("[info] Load default config %s", SectionEmailConfig)
	vars.EmailConfigSetting = new(vars.EmailConfigSettingS)
	config.MapConfig(SectionEmailConfig, vars.EmailConfigSetting)
	return nil
}
