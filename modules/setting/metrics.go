package setting

import "github.com/spf13/viper"

type MetricsService struct {
	Enabled bool
	Token   string
}

var (
	Metrics *MetricsService
)

func newMetricsService() {
	viper.SetDefault("metrics", map[string]interface{}{
		"Enabled": false,
		"Token":   "",
	})

	metricsCfg := viper.Sub("metrics")

	Metrics = &MetricsService{}
	Metrics.Enabled = metricsCfg.GetBool("Enabled")
	Metrics.Token = metricsCfg.GetString("Token")
}
