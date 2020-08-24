package sso

import (
	"strings"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Auth struct {
}

func (o *Auth) Init() error {
	return nil
}

func (o *Auth) Close() error {
	return nil
}

func (o *Auth) IsEnabled() bool {
	return true
}

func (o *Auth) userIDFormToken(c *gin.Context) uint {
	token := c.Query("token")
	if len(token) == 0 {
		token = c.Query("access_token")
	}

	if len(token) == 0 {
		header := c.GetHeader("Authorization")
		if len(header) > 0 {
			auths := strings.Fields(header)
			if len(auths) == 2 && (auths[0] == "token" || strings.ToLower(auths[0]) == "bearer") {
				token = auths[1]
			}
		}
	}

	if len(token) == 0 {
		return 0
	}

	authToken, err := models.GetAuthTokenByCode(token)
	if err != nil {
		if !models.IsErrAuthTokenEmpty(err) || !models.IsErrAuthTokenNotExist(err) {
			log.Error("Cannot get auth token by code", zap.Error(err))
		}
		return 0
	}

	if err := models.UpdateAuthToken(authToken); err != nil {
		log.Error("Cannot update auth token", zap.Error(err))
	}
	return authToken.UserID
}

func (o *Auth) VerifyAuthData(c *gin.Context, sess sessions.Session) *models.User {
	uid := o.userIDFormToken(c)
	if uid == 0 {
		return nil
	}

	user, err := models.GetUserByID(uid)
	if err != nil {
		if !models.IsErrUserNotExist(err) {
			log.Error("Cannot get user by id", zap.Error(err))
		}
		return nil
	}
	return user
}
