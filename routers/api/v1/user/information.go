package user

import (
	"net/http"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/context"
	ecode "github.com/czhj/ahfs/routers/api/v1/errcode"
)

type modifyUserInformationForm struct {
	Nickname string `json:"nickname" form:"nickname" binding:"required,gt=0,lt=16"`
}

func ModifyUserInformation(c *context.APIContext) {
	form := &modifyUserInformationForm{}
	if err := c.Bind(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, nil)
		return
	}

	user := c.User
	user.Nickname = form.Nickname

	if err := models.SaveUser(user); err != nil {
		c.InternalServerError(err)
		return
	}

	c.OK(nil)
}
