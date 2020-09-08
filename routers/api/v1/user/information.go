package user

import (
	"fmt"
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

type signInResetPasswordForm struct {
	OldPassword string `json:"old_password" form:"old_password" binding:"required,min=6,max=16"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required,min=6,max=16"`
}

func SignInUserResetPasswordPost(c *context.APIContext) {
	form := &signInResetPasswordForm{}

	if err := c.Bind(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	user := c.User
	if !user.ValidatePassword(form.OldPassword) {
		c.Error(http.StatusBadRequest, ecode.UserOldPwdIncorrect, fmt.Errorf("user old password incorrect"))
		return
	}

	user.HashPassword(form.NewPassword)
	if err := models.SaveUser(user); err != nil {
		c.InternalServerError(err)
		return
	}

	c.OK(nil)
}
