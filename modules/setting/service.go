package setting

import (
	"time"

	"github.com/spf13/viper"
)

var Service struct {
	ActiveCodeLive        time.Duration
	ResetPasswordCodeLive time.Duration
	RegisterEmailConfirm  bool
	MaxFileCapacitySize   int64
}

func newService() {
	viper.SetDefault("service", map[string]interface{}{
		"activeCodeLive":        time.Duration(15) * time.Minute,
		"resetPasswordCodeLive": time.Duration(10) * time.Minute,
		"registerEmailConfirm":  true,
		"maxFileCapacitySize":   1024 * 1024 * 512, //512M
	})

	serviceCfg := viper.Sub("service")
	Service.MaxFileCapacitySize = serviceCfg.GetInt64("maxFileCapacitySize")
	Service.ActiveCodeLive = serviceCfg.GetDuration("activeCodeLive")
	Service.ResetPasswordCodeLive = serviceCfg.GetDuration("resetPasswordCodeLive")
	Service.RegisterEmailConfirm = serviceCfg.GetBool("registerEmailConfirm")
}
