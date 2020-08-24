package file

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/context"
	"github.com/czhj/ahfs/modules/convert"
)

func UploadFile(c *context.APIContext) {
	if c.User == nil {
		c.Error(http.StatusUnauthorized, nil)
		return
	}

	filename := c.PostForm("filename")
	parentID, _ := strconv.ParseUint(c.PostForm("parent_id"), 10, 64)
	fileHeader, err := c.FormFile("upload_file")
	if err != nil {
		c.InternalServerError(err)
		return
	}

	filename = strings.TrimSpace(filename)
	if len(filename) != 0 {
		fileHeader.Filename = filename
	}

	parentFile, err := models.GetFileByID(uint(parentID), c.User.ID)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusBadRequest, err)
			return
		}
		c.InternalServerError(err)
		return
	}

	file, err := models.TryUploadFile(c.User, parentFile, *fileHeader)
	if err != nil {
		if models.IsErrFileNotDirectory(err) {
			c.Error(http.StatusBadRequest, err)
		} else if models.IsErrFileMaxSizeLimit(err) {
			c.Error(http.StatusBadRequest, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	c.OK(convert.ToFile(file))
}
