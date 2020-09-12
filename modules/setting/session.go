package setting

import (
	"github.com/czhj/ahfs/modules/log"
	"github.com/spf13/viper"
)

var (
	SessionConfig = struct {
		Provider       string
		ProviderConfig string
		CookieName     string
		CookiePath     string
		Domain         string
		MaxAge         int64
		Secure         bool
		HttpOnly       bool
	}{
		MaxAge: 86400,
	}
)

func newSessionService() {
	viper.SetDefault("session", map[string]interface{}{
		"provider":        "memory",
		"provider_config": "",
		"cookie_name":     "ahfs_lalala",
		"cookie_path":     "/",
		"secure":          false,
		"http_only":       true,
		"domain":          "",
		"maxAge":          86400,
	})

	session := viper.Sub("session")
	SessionConfig.Provider = session.GetString("provider")
	SessionConfig.ProviderConfig = session.GetString("provider_config")
	SessionConfig.CookieName = session.GetString("cookie_name")
	SessionConfig.CookiePath = session.GetString("cookie_path")
	SessionConfig.Secure = session.GetBool("secure")
	SessionConfig.Domain = session.GetString("domain")
	SessionConfig.MaxAge = session.GetInt64("maxAge")
	SessionConfig.HttpOnly = session.GetBool("http_only")

	log.Info("Session Service Enabled")
}
