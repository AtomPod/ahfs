package context

import (
	"time"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	ContextKey = "ahfs-context-key"
)

type Context struct {
	*gin.Context
	Session  sessions.Session
	Link     string
	User     *models.User
	IsSigned bool
}

func (ctx *Context) IsAdmin() bool {
	return ctx.IsSigned && ctx.User.IsAdmin
}

func (ctx *Context) HasValue(name string) bool {
	_, ok := ctx.Get(name)
	return ok
}

func Contexter() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		ctx := &Context{
			Context: c,
			Session: session,
		}
		ctx.Set("PageStartTime", time.Now())
		ctx.User, _ = auth.SignInUser(c, session)

		if ctx.User != nil {
			ctx.IsSigned = true
			ctx.Set("IsSigned", ctx.IsSigned)
			ctx.Set("UserID", ctx.User.ID)
			ctx.Set("UserNickname", ctx.User.Nickname)
			ctx.Set("UserEmail", ctx.User.IsPlain)
			ctx.Set("IsAdmin", ctx.User.IsAdmin)
		}

		c.Set(ContextKey, ctx)
		c.Next()
	}
}

func ContextWrapper(f func(c *Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, _ := c.Get(ContextKey)
		ctx, _ := value.(*Context)
		f(ctx)
	}
}
