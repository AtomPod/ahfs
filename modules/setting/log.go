package setting

import (
	"encoding/json"
	"strings"

	"github.com/czhj/ahfs/modules/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func newLogService() {

	viper.SetDefault("log", map[string]interface{}{
		"base": map[string]interface{}{
			"encoding": "console",
		},
		"loggers": map[string]interface{}{
			"console": map[string]interface{}{
				"provider": "console",
			},
		},
	})
	logCfg := viper.Sub("log")

	var config string = "{}"
	base := logCfg.GetStringMap("base")
	if len(base) != 0 {
		zapConf, err := json.Marshal(base)
		if err != nil {
			log.Warn("Cannot marshal log base config", zap.Error(err))
		} else {
			config = string(zapConf)
		}
	}

	loggers := logCfg.Sub("loggers")

	if loggers != nil {
		logger := log.NewZapTeeLogger()

		allkeys := loggers.AllKeys()
		keys := make(map[string]bool)

		for _, k := range allkeys {
			key := strings.SplitN(k, ".", 2)
			if len(key) != 0 {
				keys[key[0]] = true
			}
		}

		for i := range keys {
			key := i
			if !loggers.IsSet(key + ".provider") {
				log.Warn("Failed to add logger, provider is not set", zap.String("name", key))
				continue
			}

			provider := loggers.GetString(key + ".provider")
			loggerCfg := loggers.GetStringMap(key)

			config, err := json.Marshal(loggerCfg)
			if err != nil {
				log.Warn("Cannot marshal log config", zap.String("name", key), zap.Error(err))
				continue
			}

			if err := logger.AddLogger(key, provider, string(config)); err != nil {
				log.Warn("Failed to add logger", zap.String("name", key), zap.String("provider", provider), zap.Error(err))
			}
		}

		if err := logger.Build(config); err != nil {
			log.Error("Failed to build zap log", zap.Error(err))
		} else {
			log.SetDefault(logger)
		}
	}
}
