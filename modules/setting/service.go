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
	AvatarMaxSize         int64
}

func newService() {
	viper.SetDefault("service", map[string]interface{}{
		"active_code_live":         time.Duration(15) * time.Minute,
		"reset_password_code_live": time.Duration(10) * time.Minute,
		"register_email_confirm":   true,
		"max_file_capacity_size":   1024 * 1024 * 512, //512M
		"avatar_max_size":          1024 * 1024 * 3,
	})

	serviceCfg := viper.Sub("service")
	Service.MaxFileCapacitySize = serviceCfg.GetInt64("max_file_capacity_size")
	Service.ActiveCodeLive = serviceCfg.GetDuration("active_code_live")
	Service.ResetPasswordCodeLive = serviceCfg.GetDuration("reset_password_code_live")
	Service.RegisterEmailConfirm = serviceCfg.GetBool("register_email_confirm")
	Service.AvatarMaxSize = serviceCfg.GetInt64("avatar_max_size")
}
