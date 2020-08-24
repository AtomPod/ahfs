package sso

import (
	"strings"

	"github.com/czhj/ahfs/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type OAuth2 struct {
}

func (o *OAuth2) Init() error {
	return nil
}

func (o *OAuth2) Close() error {
	return nil
}

func (o *OAuth2) IsEnabled() bool {
	return true
}

func (o *OAuth2) userIDFormToken(c *gin.Context) uint {
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

	return 0
}

func (o *OAuth2) VerifyAuthData(c *gin.Context, sess sessions.Session) *models.User {
	return nil
}
