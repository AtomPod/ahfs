package file

import (
	"fmt"
	"net/http"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/context"
	ecode "github.com/czhj/ahfs/routers/api/v1/errcode"
)

type DownloadFileForm struct {
	FileID uint `form:"file_id" uri:"file_id" json:"file_id" binding:"required"`
}

func DownloadFile(c *context.APIContext) {
	if c.User == nil {
		c.Error(http.StatusUnauthorized, ecode.UnauthorizedError, nil)
		return
	}

	form := &DownloadFileForm{}
	if err := c.BindUri(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	userID := c.User.ID
	if c.User.IsAdmin {
		userID = 0
	}

	file, err := models.GetFileByID(form.FileID, userID)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusBadRequest, ecode.FileNotExist, err)
			return
		}
		c.InternalServerError(err)
		return
	}

	if file.IsDir() {
		c.Error(http.StatusBadRequest, ecode.FileDownloadDirError, fmt.Errorf("Cannot download a directory"))
		return
	}

	c.File(file.FileName, file.LocalPath())
}
