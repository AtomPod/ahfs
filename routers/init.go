package routers

import (
	"context"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/log"
	"github.com/czhj/ahfs/modules/setting"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func NewServices() {
	setting.NewServices()
}

func initDBEngine(ctx context.Context) (err error) {

	if err := models.NewEngine(ctx, func(e *gorm.DB) error {
		return e.AutoMigrate(&models.User{}, &models.File{}, &models.AuthToken{}).Error
	}); err != nil {
		return err
	}

	return nil
}

func GlobalInit(ctx context.Context) {

	setting.NewSetting()
	NewServices()

	if setting.ServerMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if setting.ServerMode == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	if err := initDBEngine(ctx); err != nil {
		log.Fatal("GORM engine initalization failed", zap.Error(err))
	} else {
		log.Info("GORM engine initialization success")
	}
}
