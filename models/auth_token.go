package models

import (
	"crypto/sha1"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type AuthToken struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Code              string `gorm:"unique_index;not null"`
	UserID            uint   `gorm:"index"`
	User              *User  `gorm:"-"`
	HasRecentActivity bool   `gorm:"-"`
}

func (t *AuthToken) AfterFind() {
	t.HasRecentActivity = t.UpdatedAt.Add(7 * 24 * time.Hour).After(time.Now())
}

func (t *AuthToken) LoadUser() (err error) {
	if t.User == nil {
		t.User, err = GetUserByID(t.UserID)
	}
	return
}

func CreateAuthToken(t *AuthToken) error {
	code := []byte(uuid.Must(uuid.NewRandom()).String())
	sha1Code := sha1.Sum(code)
	base64Code := base64.URLEncoding.EncodeToString(sha1Code[:])
	t.Code = base64Code
	return createAuthToken(engine, t)
}

func createAuthToken(e *gorm.DB, t *AuthToken) error {
	return e.Create(t).Error
}

func DeleteAuthTokenByUserID(userID uint) error {
	result := engine.Where("user_id=?", userID).Delete(&AuthToken{})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func DeleteAuthTokenByID(id, userID uint) error {
	result := engine.Where("id=? AND user_id = ?", id, userID).Delete(&AuthToken{})
	if err := result.Error; err != nil {
		return err
	} else if result.RowsAffected == 0 {
		return ErrAuthTokenNotExist{UserID: userID, ID: id}
	}
	return nil
}

func GetAuthTokenByCode(code string) (*AuthToken, error) {
	if code == "" {
		return nil, ErrAuthTokenEmpty{}
	}

	auth := new(AuthToken)
	err := engine.Where("code = ?", code).First(auth).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrAuthTokenNotExist{Code: code}
		}
		return nil, err
	}
	return auth, nil
}

func ListAuthTokens(uid uint, listOpts ListOptions) ([]*AuthToken, error) {
	auths := make([]*AuthToken, 0)
	query := engine.Where("user_id=?", uid)

	if listOpts.Page != 0 {
		query = listOpts.SetEnginePagination(query)
	}

	err := query.Find(&auths).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrAuthTokenNotExist{UserID: uid}
		}
		return nil, err
	}
	return auths, nil
}

func UpdateAuthToken(t *AuthToken) error {
	return updateAuthToken(engine, t)
}

func updateAuthToken(e *gorm.DB, t *AuthToken) error {
	db := e.Save(t)
	if err := db.Error; err != nil {
		return err
	}

	if db.RowsAffected == 0 {
		return ErrAuthTokenNotExist{Code: t.Code, UserID: t.UserID}
	}

	return nil
}
