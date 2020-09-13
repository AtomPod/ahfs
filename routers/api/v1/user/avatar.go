package user

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/czhj/ahfs/modules/context"
	"github.com/czhj/ahfs/modules/setting"
	ecode "github.com/czhj/ahfs/routers/api/v1/errcode"
)

type updateAvatarForm struct {
	Avatar *multipart.FileHeader `form:"avatar" binding:"required"`
}

func UpdateAvatar(c *context.APIContext) {
	form := &updateAvatarForm{}

	if err := c.ShouldBind(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	user := c.User
	if form.Avatar.Size > setting.Service.AvatarMaxSize {
		c.Error(http.StatusBadRequest, ecode.UserAvatarSizeTooLarge, fmt.Errorf("avatar file size too large"))
		return
	}

	avatar, err := form.Avatar.Open()
	if err != nil {
		c.InternalServerError(err)
		return
	}
	defer avatar.Close()

	avatarData, err := ioutil.ReadAll(avatar)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	if err := user.UploadAvatar(avatarData); err != nil {
		c.InternalServerError(err)
		return
	}

	c.OK(nil)
}
