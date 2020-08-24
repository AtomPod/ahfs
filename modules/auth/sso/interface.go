package sso

import (
	"github.com/czhj/ahfs/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SingleSignOn interface {
	Init() error
	Close() error
	IsEnabled() bool
	VerifyAuthData(ctx *gin.Context, sess sessions.Session) *models.User
}
