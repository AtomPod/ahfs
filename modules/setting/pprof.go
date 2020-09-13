package setting

import "github.com/spf13/viper"

type Pprof struct {
	Enabled  bool
	HTTPAddr string
	HTTPPort string
}

var (
	PprofService *Pprof
)

func newPprofService() {
	viper.SetDefault("pprof", map[string]interface{}{
		"enabled":   false,
		"http_addr": "0.0.0.0",
		"http_port": "6060",
	})

	pprofCfg := viper.Sub("pprof")
	PprofService = new(Pprof)
	PprofService.Enabled = pprofCfg.GetBool("enabled")
	PprofService.HTTPAddr = pprofCfg.GetString("http_addr")
	PprofService.HTTPPort = pprofCfg.GetString("http_port")
}
