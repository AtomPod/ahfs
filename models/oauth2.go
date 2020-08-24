package models

import (
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/jinzhu/gorm"
)

var (
	zeroUTFTime = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
)

type OAuth2Token struct {
	ID               uint   `gorm:"primary_key"`
	ClientID         string `gorm:"index"`
	UserID           string `gorm:"index"`
	RedirectURI      string
	Scope            string
	Code             string `gorm:"index"`
	CodeExpiredAt    time.Time
	CodeExpiresIn    time.Duration
	Access           string `gorm:"index"`
	AccessExpiredAt  time.Time
	AccessExpiresIn  time.Duration
	Refresh          string `gorm:"index"`
	RefreshExpiredAt time.Time
	RefreshExpiresIn time.Duration
}

func (t *OAuth2Token) New() oauth2.TokenInfo {
	token := new(OAuth2Token)
	return token
}

func (t *OAuth2Token) GetClientID() string {
	return t.ClientID
}

func (t *OAuth2Token) SetClientID(id string) {
	t.ClientID = id
}

func (t *OAuth2Token) GetUserID() string {
	return t.UserID
}

func (t *OAuth2Token) SetUserID(id string) {
	t.UserID = id
}

func (t *OAuth2Token) GetRedirectURI() string {
	return t.RedirectURI
}

func (t *OAuth2Token) SetRedirectURI(uri string) {
	t.RedirectURI = uri
}

func (t *OAuth2Token) GetScope() string {
	return t.Scope
}

func (t *OAuth2Token) SetScope(scope string) {
	t.Scope = scope
}

func (t *OAuth2Token) GetCode() string {
	return t.Code
}

func (t *OAuth2Token) SetCode(code string) {
	t.Code = code
}

func (t *OAuth2Token) GetCodeCreateAt() time.Time {
	return t.CodeExpiredAt.Add(-t.CodeExpiresIn)
}

func (t *OAuth2Token) SetCodeCreateAt(createAt time.Time) {
	t.CodeExpiredAt = createAt.Add(t.CodeExpiresIn)
}

func (t *OAuth2Token) GetCodeExpiresIn() time.Duration {
	return t.CodeExpiresIn
}

func (t *OAuth2Token) SetCodeExpiresIn(d time.Duration) {
	t.CodeExpiredAt = t.CodeExpiredAt.Add(d - t.CodeExpiresIn)
	t.CodeExpiresIn = d
}

func (t *OAuth2Token) GetAccess() string {
	return t.Access
}

func (t *OAuth2Token) SetAccess(token string) {
	t.Access = token
}

func (t *OAuth2Token) GetAccessCreateAt() time.Time {
	return t.AccessExpiredAt.Add(-t.AccessExpiresIn)
}

func (t *OAuth2Token) SetAccessCreateAt(createAt time.Time) {
	t.AccessExpiredAt = createAt.Add(t.AccessExpiresIn)
}

func (t *OAuth2Token) GetAccessExpiresIn() time.Duration {
	return t.AccessExpiresIn
}

func (t *OAuth2Token) SetAccessExpiresIn(d time.Duration) {
	t.AccessExpiredAt = t.AccessExpiredAt.Add(d - t.AccessExpiresIn)
	t.AccessExpiresIn = d
}

func (t *OAuth2Token) GetRefresh() string {
	return t.Refresh
}

func (t *OAuth2Token) SetRefresh(token string) {
	t.Refresh = token
}

func (t *OAuth2Token) GetRefreshCreateAt() time.Time {
	return t.RefreshExpiredAt.Add(-t.RefreshExpiresIn)
}

func (t *OAuth2Token) SetRefreshCreateAt(createAt time.Time) {
	t.RefreshExpiredAt = createAt.Add(t.RefreshExpiresIn)
}

func (t *OAuth2Token) GetRefreshExpiresIn() time.Duration {
	return t.RefreshExpiresIn
}

func (t *OAuth2Token) SetRefreshExpiresIn(d time.Duration) {
	t.RefreshExpiredAt = t.RefreshExpiredAt.Add(d - t.RefreshExpiresIn)
	t.RefreshExpiresIn = d
}

func (t *OAuth2Token) Invalid() error {
	return t.invalid(engine)
}

func (t *OAuth2Token) invalid(e *gorm.DB) error {
	return e.Delete(t).Error
}

func (t *OAuth2Token) InvalidCode() error {
	return t.invalidCode(engine)
}

func (t *OAuth2Token) invalidCode(e *gorm.DB) error {
	t.CodeExpiredAt = zeroUTFTime
	return e.Model(t).Update("code_expired_at", t.CodeExpiredAt).Error
}

func (t *OAuth2Token) InvalidAccessToken() error {
	return t.invalidAccessToken(engine)
}

func (t *OAuth2Token) invalidAccessToken(e *gorm.DB) error {
	t.AccessExpiredAt = zeroUTFTime
	return e.Model(t).Update("access_expired_at", t.AccessExpiredAt).Error
}

func (t *OAuth2Token) InvalidRefreshToken() error {
	return t.invalidRefreshToken(engine)
}

func (t *OAuth2Token) invalidRefreshToken(e *gorm.DB) error {
	t.RefreshExpiredAt = zeroUTFTime
	return e.Model(t).Update("refresh_expired_at", t.RefreshExpiredAt).Error
}

func CreateOAuth2Token(t *OAuth2Token) error {
	return createOAuth2Token(engine, t)
}

func createOAuth2Token(e *gorm.DB, t *OAuth2Token) error {
	return e.Create(t).Error
}

func GetOAuth2TokenByCode(code string) (*OAuth2Token, error) {
	return getOAuth2TokenByCode(engine, code)
}

func getOAuth2TokenByCode(e *gorm.DB, code string) (*OAuth2Token, error) {
	token := new(OAuth2Token)
	err := e.Where("code = ? AND code_expired_at > ?", code, time.Now()).
		First(token).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrOAuth2TokenNotExist{Code: code}
		}
		return nil, err
	}
	return token, nil
}

func GetOAuth2TokenByAccessToken(token string) (*OAuth2Token, error) {
	return getOAuth2TokenByAccessToken(engine, token)
}

func getOAuth2TokenByAccessToken(e *gorm.DB, accessToken string) (*OAuth2Token, error) {
	token := new(OAuth2Token)
	err := e.Where("access = ? AND access_expired_at > ?", accessToken, time.Now()).
		First(token).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrOAuth2TokenNotExist{AccessToken: accessToken}
		}
		return nil, err
	}
	return token, nil
}

func GetOAuth2TokenByRefreshToken(token string) (*OAuth2Token, error) {
	return getOAuth2TokenByRefreshToken(engine, token)
}

func getOAuth2TokenByRefreshToken(e *gorm.DB, refresh string) (*OAuth2Token, error) {
	token := new(OAuth2Token)
	err := e.Where("refresh = ? AND refresh_expired_at > ?", refresh, time.Now()).
		First(token).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrOAuth2TokenNotExist{RefreshToken: refresh}
		}
		return nil, err
	}
	return token, nil
}

func DeleteOAuth2TokenByCode(code string) error {
	return deleteOAuth2TokenByCode(engine, code)
}

func deleteOAuth2TokenByCode(e *gorm.DB, code string) error {
	db := e.Where("code=? AND code_expired_at > ?", code, time.Now()).
		Update("code_expired_at", zeroUTFTime)
	if err := db.Error; err != nil {
		return err
	}
	if db.RowsAffected == 0 {
		return ErrOAuth2TokenNotExist{Code: code}
	}
	return nil
}

func DeleteOAuth2TokenByAccessToken(accessToken string) error {
	return deleteOAuth2TokenByAccessToken(engine, accessToken)
}

func deleteOAuth2TokenByAccessToken(e *gorm.DB, accessToken string) error {
	db := e.Where("access=? AND access_expired_at > ?", accessToken, time.Now()).
		Update("access_expired_at", zeroUTFTime)
	if err := db.Error; err != nil {
		return err
	}
	if db.RowsAffected == 0 {
		return ErrOAuth2TokenNotExist{AccessToken: accessToken}
	}
	return nil
}

func DeleteOAuth2TokenByRefreshToken(refreshToken string) error {
	return deleteOAuth2TokenByRefreshToken(engine, refreshToken)
}

func deleteOAuth2TokenByRefreshToken(e *gorm.DB, refreshToken string) error {
	db := e.Where("refresh=? AND refresh_expired_at > ?", refreshToken, time.Now()).
		Update("refresh_expired_at", zeroUTFTime)
	if err := db.Error; err != nil {
		return err
	}
	if db.RowsAffected == 0 {
		return ErrOAuth2TokenNotExist{RefreshToken: refreshToken}
	}
	return nil
}
