package setting

import (
	"os"
	"path/filepath"
	"time"

	"github.com/czhj/ahfs/modules/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type DatabaseConfig struct {
	Driver       string
	URL          string
	MaxIdleConns int
	MaxOpenConns int
	MaxListTime  time.Duration
}

var (
	EnableSQLite3 bool
	Database      DatabaseConfig
)

func newDBService() {
	viper.SetDefault("database", map[string]interface{}{
		"driver":       "sqlite3",
		"url":          "/db/gorm.db",
		"maxIdleConns": 0,
		"maxOpenConns": 0,
		"maxLifeTime":  "0s",
	})

	dbcfg := viper.Sub("database")

	Database.Driver = dbcfg.GetString("driver")
	Database.URL = dbcfg.GetString("url")
	Database.MaxIdleConns = dbcfg.GetInt("maxIdleConns")
	Database.MaxOpenConns = dbcfg.GetInt("maxOpenConns")
	Database.MaxListTime = dbcfg.GetDuration("maxListTime")

	if Database.Driver == "sqlite3" {
		EnableSQLite3 = true

		if !filepath.IsAbs(Database.URL) {
			Database.URL = filepath.Join(AppDataPath, Database.URL)
		}

		dir := filepath.Dir(Database.URL)
		if err := os.MkdirAll(dir, os.ModeDir); err != nil {
			log.Error("could not create directory", zap.String("path", dir), zap.Error(err))
		}
	}
}
