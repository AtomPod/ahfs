package user

import (
	"fmt"
	"net/http"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/cache"
	"github.com/czhj/ahfs/modules/code"
	"github.com/czhj/ahfs/modules/context"
	"github.com/czhj/ahfs/modules/convert"
	ecode "github.com/czhj/ahfs/routers/api/v1/errcode"
	"github.com/czhj/ahfs/services/mailer"
)

type RequestResetPwdCodeForm struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}

func RequestResetPwdCode(c *context.APIContext) {
	form := &RequestResetPwdCodeForm{}

	if err := c.Bind(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.EmailFormatError, err)
		return
	}

	if used, err := models.IsEmailUsed(form.Email); err != nil || !used {
		if err != nil {
			c.InternalServerError(err)
			return
		}
		c.Error(http.StatusNotFound, ecode.EmailNotFound, fmt.Errorf("email [%s] does not exists", form.Email))
		return
	}

	resetCode, err := code.CreateEmailResetPwdCode(form.Email)
	if err == code.ErrTooOften {
		c.Error(http.StatusBadRequest, ecode.EmailResetPwdCodeTooOften, err)
		return
	}

	if err != nil {
		c.InternalServerError(err)
		return
	}

	mailer.SendResetPwdCodeMail(form.Email, resetCode)
	c.OK(nil)
}

type ResetPwdCodeForm struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Code     string `json:"code" form:"code" binding:"required"`
	Password string `json:"password" form:"password" binding:"required,min=6,max=16"`
}

func ResetPasswordPost(c *context.APIContext) {
	form := &ResetPwdCodeForm{}

	if err := c.Bind(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	ok, err := code.VerifyEmailResetPwdCode(form.Email, form.Code)
	if err != nil || !ok {
		if !ok || err == cache.ErrNotFound {
			c.Error(http.StatusNotFound, ecode.EmailResetPwdCodeError, fmt.Errorf("reset password verification code does not exists"))
		} else {
			c.InternalServerError(err)
		}
		return
	}

	user, err := models.GetUserByEmail(form.Email)
	if err != nil {
		if models.IsErrUserNotExist(err) {
			c.Error(http.StatusNotFound, ecode.EmailNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	user.HashPassword(form.Password)
	if err := models.SaveUser(user); err != nil {
		c.InternalServerError(err)
		return
	}

	code.RemoveEmailResetPwdCode(form.Email, form.Code)

	if err := models.DeleteAuthTokenByUserID(user.ID); err != nil {
		c.InternalServerError(err)
		return
	}

	authToken := &models.AuthToken{UserID: user.ID}
	if err := models.CreateAuthToken(authToken); err != nil {
		c.InternalServerError(fmt.Errorf("CreateAuthToken: Unable to create auth token: %v", err))
		return
	}

	result := convert.ToToken(authToken)
	c.OK(result)
}

type EditUserPasswordForm struct {
	OldPassword string `json:"old_password" form:"old_password" binding:"required,min=6,max=16"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required,min=6,max=16"`
}

func EditUserPassword(c *context.APIContext) {
	form := &EditUserPasswordForm{}

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
