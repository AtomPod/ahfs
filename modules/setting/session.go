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
		"provider":       "memory",
		"providerConfig": "",
		"cookieName":     "ahfs_lalala",
		"cookiePath":     "/",
		"secure":         false,
		"httpOnly":       true,
		"domain":         "",
		"maxAge":         86400,
	})

	session := viper.Sub("session")
	SessionConfig.Provider = session.GetString("provider")
	SessionConfig.ProviderConfig = session.GetString("providerConfig")
	SessionConfig.CookieName = session.GetString("cookieName")
	SessionConfig.CookiePath = session.GetString("cookiePath")
	SessionConfig.Secure = session.GetBool("secure")
	SessionConfig.Domain = session.GetString("domain")
	SessionConfig.MaxAge = session.GetInt64("maxAge")
	SessionConfig.HttpOnly = session.GetBool("httpOnly")

	log.Info("Session Service Enabled")
}
