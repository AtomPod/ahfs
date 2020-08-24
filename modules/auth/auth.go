package auth

import (
	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/auth/sso"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SignInUser(ctx *gin.Context, sess sessions.Session) (*models.User, bool) {
	for _, method := range sso.Methods() {
		if !method.IsEnabled() {
			continue
		}

		user := method.VerifyAuthData(ctx, sess)
		if user != nil {
			return user, true
		}
	}
	return nil, false
}
