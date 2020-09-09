package user

import (
	"net/http"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/context"
	ecode "github.com/czhj/ahfs/routers/api/v1/errcode"
)

type EditUserOption struct {
	Nickname *string `json:"nickname" form:"nickname" binding:"omitempty,gt=0,lt=16"`
}

func EditUser(c *context.APIContext) {
	form := &EditUserOption{}
	if err := c.Bind(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, nil)
		return
	}

	user := c.User

	if form.Nickname != nil {
		user.Nickname = *form.Nickname
	}

	if err := models.SaveUser(user); err != nil {
		c.InternalServerError(err)
		return
	}

	c.OK(nil)
}
