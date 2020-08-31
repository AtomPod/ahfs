package user

import (
	"fmt"
	"net/http"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/code"
	"github.com/czhj/ahfs/modules/context"
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
	if err != nil {
		c.InternalServerError(err)
		return
	}

	mailer.SendResetPwdCodeMail(form.Email, resetCode)
	c.OK(nil)
}
