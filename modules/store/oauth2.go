package store

import (
	"context"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/convert"
	"github.com/go-oauth2/oauth2/v4"
)

type OAuth2TokenStore struct {
}

func (store *OAuth2TokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	token := convert.FromOAuth2TokenInfo(info)
	return models.CreateOAuth2Token(token)
}

// delete the authorization code
func (store *OAuth2TokenStore) RemoveByCode(ctx context.Context, code string) error {
	return models.DeleteOAuth2TokenByCode(code)
}

// use the access token to delete the token information
func (store *OAuth2TokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return models.DeleteOAuth2TokenByAccessToken(access)
}

// use the refresh token to delete the token information
func (store *OAuth2TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return models.DeleteOAuth2TokenByRefreshToken(refresh)
}

// use the authorization code for token information data
func (store *OAuth2TokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	return models.GetOAuth2TokenByCode(code)
}

// use the access token for token information data
func (store *OAuth2TokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	return models.GetOAuth2TokenByAccessToken(access)
}

// use the refresh token for token information data
func (store *OAuth2TokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	return models.GetOAuth2TokenByRefreshToken(refresh)
}
