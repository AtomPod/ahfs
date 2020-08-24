package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type OAuth2Application struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	UID  uint  `gorm:"index"`
	User *User `gorm:"-"`

	Name         string
	ClientID     string `gorm:"unique_index"`
	ClientSecret string
	Domain       string
}

func (app *OAuth2Application) LoadUser() (err error) {
	app.User, err = GetUserByID(app.UID)
	return
}

func CreateOAuth2Application(app *OAuth2Application) error {
	return createOAuth2Application(engine, app)
}

func createOAuth2Application(e *gorm.DB, app *OAuth2Application) error {
	return e.Create(app).Error
}

func DeleteOAuth2Application(app *OAuth2Application) error {
	return deleteOAuth2Application(engine, app)
}

func deleteOAuth2Application(e *gorm.DB, app *OAuth2Application) error {
	db := e.Where(app).Delete(&OAuth2Application{})
	if err := db.Error; err != nil {
		return err
	}

	if db.RowsAffected == 0 {
		return ErrOAuth2ApplicationNotExist{ID: app.ID, Name: app.Name, ClientID: app.ClientID, UserID: app.UID}
	}
	return nil
}

func GetOAuth2ApplicationByUserID(userID uint) ([]*OAuth2Application, error) {
	return getOAuth2ApplicationByUserID(engine, userID)
}

func getOAuth2ApplicationByUserID(e *gorm.DB, userID uint) ([]*OAuth2Application, error) {
	apps := make([]*OAuth2Application, 0)
	db := e.Where("uid=?", userID).Find(&apps)
	if err := db.Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrOAuth2ApplicationNotExist{UserID: userID}
		}
		return nil, err
	}
	return apps, nil
}

func GetOAuth2ApplicationByClientID(clientID string) (*OAuth2Application, error) {
	return getOAuth2ApplicationByClientID(engine, clientID)
}

func getOAuth2ApplicationByClientID(e *gorm.DB, clientID string) (*OAuth2Application, error) {
	app := new(OAuth2Application)
	db := e.Where("client_id=?", clientID).First(app)
	if err := db.Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrOAuth2ApplicationNotExist{ClientID: clientID}
		}
		return nil, err
	}
	return app, nil
}
