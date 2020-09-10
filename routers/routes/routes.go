package routes

import (
	"path/filepath"
	"time"

	"github.com/czhj/ahfs/modules/metrics"
	"github.com/czhj/ahfs/modules/templates"
	"github.com/czhj/ahfs/routers"
	"github.com/gin-contrib/gzip"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/czhj/ahfs/modules/context"
	"github.com/czhj/ahfs/modules/log"
	"github.com/czhj/ahfs/modules/session"
	"github.com/czhj/ahfs/modules/setting"
	v1 "github.com/czhj/ahfs/routers/api/v1"
	"github.com/czhj/ahfs/services/mailer"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func NewEngine() *gin.Engine {
	engine := gin.New()

	engine.Use(ginzap.Ginzap(log.ZapLogger(), time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(log.ZapLogger(), true))

	engine.Static("/static", filepath.Join(setting.StaticRootPath, "public"))
	engine.Static("/avatars", setting.AvatarUploadPath)

	if setting.EnableGzip {
		engine.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	engine.Use(session.NewSession(session.Options{
		Provider:       setting.SessionConfig.Provider,
		ProviderConfig: setting.SessionConfig.ProviderConfig,
		HttpOnly:       setting.SessionConfig.HttpOnly,
		Secure:         setting.SessionConfig.Secure,
		CookieName:     setting.SessionConfig.CookieName,
		CookiePath:     setting.SessionConfig.CookiePath,
		MaxAge:         setting.SessionConfig.MaxAge,
		Domain:         setting.SessionConfig.Domain,
		Secret:         []byte("test session"),
	}))

	engine.Use(context.Contexter())
	mailer.InitMailTemplate(templates.Mailer())

	if setting.Metrics.Enabled {
		collector := metrics.NewCollector()
		prometheus.MustRegister(collector)

		engine.GET("/metrics", routers.Metrics)
	}

	return engine
}

func RegisterRoutes(e *gin.Engine) {
	api := e.Group("/api", context.APIContexter())
	v1.RegisterRoutes(api)
}
