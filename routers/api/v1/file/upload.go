package file

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/context"
	"github.com/czhj/ahfs/modules/convert"
	"github.com/czhj/ahfs/modules/validator"
	ecode "github.com/czhj/ahfs/routers/api/v1/errcode"
)

func UploadFile(c *context.APIContext) {

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

	if !validator.ValidFilename(fileHeader.Filename) {
		c.Error(http.StatusBadRequest, ecode.FilenameFormatError, fmt.Errorf(`filename couldn't containt \/:*?"<>|`))
		return
	}

	var parentFile *models.File
	if parentID == 0 {
		parentFile, err = models.GetUserRootFile(c.User.ID)
	} else {
		parentFile, err = models.GetFileByID(uint(parentID), c.User.ID)
	}

	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusBadRequest, ecode.FileNotExist, err)
			return
		}
		c.InternalServerError(err)
		return
	}

	file, err := models.TryUploadFile(c.User, parentFile, fileHeader)
	if err != nil {
		if models.IsErrFileNotDirectory(err) {
			c.Error(http.StatusBadRequest, ecode.FileNotDirError, err)
		} else if models.IsErrFileMaxSizeLimit(err) {
			c.Error(http.StatusBadRequest, ecode.FileStorageFulled, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	c.OK(convert.ToFile(file))
}
