package file

import (
	"net/http"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/convert"
	api "github.com/czhj/ahfs/modules/structs"

	"github.com/czhj/ahfs/modules/context"
)

type ReadDirectoryForm struct {
	UserID uint `form:"user_id" uri:"user_id" json:"user_id" binding:"omitempty"`
	FileID uint `form:"file_id" uri:"file_id" json:"file_id" binding:"required"`
}

func ReadDirectory(c *context.APIContext) {
	if c.User == nil {
		c.Error(http.StatusUnauthorized, nil)
		return
	}

	form := &ReadDirectoryForm{}
	if err := c.Bind(form); err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	userID := c.User.ID
	if c.IsAdmin() {
		userID = form.UserID
	}

	directory, err := models.GetFileByID(form.FileID, userID)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	files, err := directory.ReadDir()
	if err != nil {
		if models.IsErrFileNotDirectory(err) {
			c.Error(http.StatusBadRequest, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	outfiles := make([]*api.File, len(files))
	for i, file := range files {
		outfiles[i] = convert.ToFile(file)
	}
	c.OK(outfiles)
}
