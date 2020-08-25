package file

import (
	"net/http"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/convert"
	api "github.com/czhj/ahfs/modules/structs"

	"github.com/czhj/ahfs/modules/context"
)

type ReadDirectoryUriForm struct {
	FileID uint `form:"file_id" uri:"file_id" json:"file_id" binding:"required"`
}

type ReadDirectoryForm struct {
	UserID uint `form:"user_id" json:"user_id" binding:"omitempty"`
}

func ReadDirectory(c *context.APIContext) {
	if c.User == nil {
		c.Error(http.StatusUnauthorized, nil)
		return
	}

	uriForm := &ReadDirectoryUriForm{}
	if err := c.BindUri(uriForm); err != nil {
		c.Error(http.StatusBadRequest, err)
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

	directory, err := models.GetFileByID(uriForm.FileID, userID)
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

type UserRootDirForm struct {
	// 只有管理员可以使用
	UserID uint `form:"user_id" query:"user_id" uri:"user_id" json:"user_id" binding:"omitempty"`
}

func GetUserRootDirectory(c *context.APIContext) {
	if !c.IsSigned {
		c.Error(http.StatusUnauthorized, nil)
		return
	}

	form := &UserRootDirForm{}

	if c.IsAdmin() {
		if err := c.Bind(form); err != nil {
			c.Error(http.StatusBadRequest, err)
			return
		}
	}

	userID := c.User.ID
	if c.IsAdmin() && form.UserID != 0 {
		userID = form.UserID
	}

	root, err := models.GetUserRootFile(userID)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	files, err := root.ReadDir()
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	apiFiles := make([]*api.File, len(files))
	for i, file := range files {
		apiFiles[i] = convert.ToFile(file)
	}
	c.OK(apiFiles)
}

type RenameFileUriForm struct {
	FileID uint `form:"file_id" uri:"file_id"  json:"file_id" binding:"required"`
}

type RenameFileForm struct {
	FileName string `form:"filename" json:"filename" binding:"required,min=6,max=32"`
}

func RenameFile(c *context.APIContext) {
	if !c.IsSigned {
		c.Error(http.StatusUnauthorized, nil)
		return
	}

	uriform := &RenameFileUriForm{}
	if err := c.BindUri(uriform); err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	form := &RenameFileForm{}
	if err := c.Bind(form); err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	userID := c.User.ID
	if c.IsAdmin() {
		userID = 0
	}

	file, err := models.GetFileByID(uriform.FileID, userID)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	file.FileName = form.FileName
	if err := models.RenameFile(file); err != nil {
		if models.IsErrModifyRootFile(err) {
			c.Error(http.StatusBadRequest, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	c.OK(convert.ToFile(file))
}

type DeleteFileForm struct {
	FileID uint `form:"file_id" uri:"file_id" json:"file_id" binding:"required"`
}

func DeleteFile(c *context.APIContext) {
	if !c.IsSigned {
		c.Error(http.StatusUnauthorized, nil)
		return
	}

	form := &DeleteFileForm{}
	if err := c.BindUri(form); err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	userID := c.User.ID
	if c.IsAdmin() {
		userID = 0
	}

	file, err := models.GetFileByID(form.FileID, userID)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	if err := models.DeleteFile(file); err != nil {
		c.InternalServerError(err)
		return
	}

	c.OK(nil)
}

type MoveFileUriForm struct {
	FileID uint `form:"file_id" uri:"file_id" json:"file_id" binding:"required"`
}

type MoveFileForm struct {
	DirectoryID uint `form:"directory_id" json:"directory_id" binding:"required"`
}

func MoveFile(c *context.APIContext) {
	if !c.IsSigned {
		c.Error(http.StatusUnauthorized, nil)
		return
	}

	uriform := &MoveFileUriForm{}
	if err := c.BindUri(uriform); err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	form := &MoveFileForm{}
	if err := c.Bind(form); err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	userID := c.User.ID
	if c.IsAdmin() {
		userID = 0
	}

	file, err := models.GetFileByID(uriform.FileID, userID)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	diretcory, err := models.GetFileByID(form.DirectoryID, userID)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	if err := models.MoveFile(file, diretcory); err != nil {
		if models.IsErrModifyRootFile(err) {
			c.Error(http.StatusBadRequest, err)
		} else if models.IsErrFileParentNotDirectory(err) {
			c.Error(http.StatusBadRequest, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	c.OK(nil)
}
