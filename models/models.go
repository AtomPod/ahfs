package models

import (
	"context"
	"fmt"

	"github.com/czhj/ahfs/modules/setting"
	"github.com/jinzhu/gorm"
)

var (
	engine *gorm.DB
)

func getEngine() (*gorm.DB, error) {
	driver := setting.Database.Driver
	url := setting.Database.URL

	db, err := gorm.Open(driver, url)
	if err != nil {
		return nil, err
	}
	if setting.ServerMode != "release" {
		db = db.Debug()
	}
	db.DB().SetConnMaxLifetime(setting.Database.MaxListTime)
	db.DB().SetMaxIdleConns(setting.Database.MaxIdleConns)
	db.DB().SetMaxOpenConns(setting.Database.MaxOpenConns)
	return db, nil
}

func NewEngine(ctx context.Context, migrateFunc func(*gorm.DB) error) error {
	db, err := getEngine()
	if err != nil {
		return fmt.Errorf("Failed to connect sql: %v", err)
	}
	engine = db

	if err := engine.DB().PingContext(ctx); err != nil {
		return fmt.Errorf("Failed to ping sql: %v", err)
	}

	if err := migrateFunc(engine); err != nil {
		return fmt.Errorf("migrate: %v", err)
	}

	return nil
}
