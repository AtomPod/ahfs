package store

import (
	"context"
	"strconv"

	"github.com/czhj/ahfs/models"
	"github.com/go-oauth2/oauth2/v4"
	oauth2Models "github.com/go-oauth2/oauth2/v4/models"
)

type OAuth2ApplicationStore struct {
}

func (store *OAuth2ApplicationStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	app, err := models.GetOAuth2ApplicationByClientID(id)
	if err != nil {
		return nil, err
	}
	client := &oauth2Models.Client{
		ID:     app.ClientID,
		Secret: app.ClientSecret,
		Domain: app.Domain,
		UserID: strconv.FormatUint(uint64(app.UID), 10),
	}
	return client, nil
}
