package users

import (
	"net/http"

	"github.com/czhj/ahfs/modules/auth"
	"github.com/czhj/ahfs/modules/context"
)

func SignInPost(ctx *context.Context) {
	if ctx.IsSigned {
		ctx.Redirect(http.StatusMovedPermanently, "/")
		return
	}

	var form auth.SignInForm
	if err := ctx.Bind(&form); err != nil {

	}
}
