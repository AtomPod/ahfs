package sso

import (
	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Session struct {
}

func (s *Session) Init() error {
	return nil
}

func (s *Session) Close() error {
	return nil
}

func (s *Session) IsEnabled() bool {
	return true
}

func (s *Session) VerifyAuthData(ctx *gin.Context, sess sessions.Session) *models.User {

	uid := sess.Get("uid")
	if uid == nil {
		return nil
	}

	id, ok := uid.(uint)
	if !ok {
		return nil
	}

	user, err := models.GetUserByID(id)
	if err != nil {
		if !models.IsErrUserNotExist(err) {
			log.Error("GetUserByID", zap.Error(err))
		}
		return nil
	}

	return user
}
