package file

import (
	"net/http"
	"strconv"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/convert"
	api "github.com/czhj/ahfs/modules/structs"
	ecode "github.com/czhj/ahfs/routers/api/v1/errcode"

	"github.com/czhj/ahfs/modules/context"
)

// type ReadDirectoryUriForm struct {
// 	FileID uint `form:"file_id" uri:"file_id" json:"file_id" binding:"required"`
// }

type ReadDirectoryForm struct {
	UserID  uint `form:"user_id" query:"user_id" binding:"omitempty"`
	OnlyDir bool `form:"only_dir" query:"only_dir" binding:"omitempty"`
}

func ReadDirectory(c *context.APIContext) {

	// uriForm := &ReadDirectoryUriForm{}
	// if err := c.BindUri(uriForm); err != nil {
	// 	c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
	// 	return
	// }

	fileIDParam := c.Param("file_id")
	fileID, err := strconv.ParseUint(fileIDParam, 10, 64)
	if err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	form := &ReadDirectoryForm{}
	if err := c.ShouldBind(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	userID := c.User.ID
	if c.IsAdmin() {
		userID = form.UserID
	}

	directory, err := models.GetFileByID(uint(fileID), userID)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusNotFound, ecode.FileNotExist, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	files, err := directory.ReadDir(models.ReadDirOption{
		OnlyDir: form.OnlyDir,
	})
	if err != nil {
		if models.IsErrFileNotDirectory(err) {
			c.Error(http.StatusBadRequest, ecode.FileNotDirError, err)
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
	form := &UserRootDirForm{}

	if c.IsAdmin() {
		if err := c.ShouldBind(form); err != nil {
			c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
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
			c.Error(http.StatusNotFound, ecode.FileNotExist, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	files, err := root.ReadDir(models.ReadDirOption{})
	if err != nil {
		if models.IsErrFileNotDirectory(err) {
			c.Error(http.StatusNotFound, ecode.FileNotDirError, err)
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

// type RenameFileUriForm struct {
// 	FileID uint `form:"file_id" uri:"file_id"  json:"file_id" binding:"required"`
// }

type RenameFileForm struct {
	FileName string `form:"filename" json:"filename" binding:"required,min=6,max=32"`
}

func RenameFile(c *context.APIContext) {

	// uriform := &RenameFileUriForm{}
	// if err := c.BindUri(uriform); err != nil {
	// 	c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
	// 	return
	// }

	// }
	fileIDParam := c.Param("file_id")
	fileID, err := strconv.ParseUint(fileIDParam, 10, 64)
	if err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	form := &RenameFileForm{}
	if err := c.ShouldBind(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	userID := c.User.ID
	if c.IsAdmin() {
		userID = 0
	}

	file, err := models.GetFileByID(uint(fileID), userID)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusNotFound, ecode.FileNotExist, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	file.FileName = form.FileName
	if err := models.RenameFile(file); err != nil {
		if models.IsErrModifyRootFile(err) {
			c.Error(http.StatusBadRequest, ecode.FileRootOperateError, err)
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

	// form := &DeleteFileForm{}
	// if err := c.BindUri(form); err != nil {
	// 	c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
	// 	return
	// }
	fileIDParam := c.Param("file_id")
	fileID, err := strconv.ParseUint(fileIDParam, 10, 64)
	if err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	userID := c.User.ID
	if c.IsAdmin() {
		userID = 0
	}

	file, err := models.GetFileByID(uint(fileID), userID)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusNotFound, ecode.FileNotExist, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	if err := models.DeleteFile(file); err != nil {
		if models.IsErrModifyRootFile(err) {
			c.Error(http.StatusBadRequest, ecode.FileRootOperateError, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	c.OK(nil)
}

// type MoveFileUriForm struct {
// 	FileID uint `form:"file_id" uri:"file_id" json:"file_id" binding:"required"`
// }

type MoveFileForm struct {
	DirectoryID uint `form:"directory_id" json:"directory_id" binding:"required"`
}

func MoveFile(c *context.APIContext) {

	// uriform := &MoveFileUriForm{}
	// if err := c.BindUri(uriform); err != nil {
	// 	c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
	// 	return
	// }

	fileIDParam := c.Param("file_id")
	fileID, err := strconv.ParseUint(fileIDParam, 10, 64)
	if err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	form := &MoveFileForm{}
	if err := c.ShouldBind(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	userID := c.User.ID
	if c.IsAdmin() {
		userID = 0
	}

	file, err := models.GetFileByID(uint(fileID), userID)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusNotFound, ecode.FileNotExist, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	diretcory, err := models.GetFileByID(form.DirectoryID, userID)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusNotFound, ecode.FileDirNotExists, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	if err := models.MoveFile(file, diretcory); err != nil {
		if models.IsErrModifyRootFile(err) {
			c.Error(http.StatusBadRequest, ecode.FileRootOperateError, err)
		} else if models.IsErrFileParentNotDirectory(err) {
			c.Error(http.StatusBadRequest, ecode.FileParentNotDirError, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	c.OK(nil)
}

type CreateDirForm struct {
	ParentID      uint   `json:"parent_id" form:"parent_id" binding:"omitempty"`
	DirectoryName string `json:"directory_name" form:"directory_name" binding:"required,gt=0,lt=256"`
}

func CreateDirectory(c *context.APIContext) {

	form := &CreateDirForm{}
	if err := c.ShouldBind(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	var dir *models.File
	var err error
	if form.ParentID == 0 {
		dir, err = models.GetUserRootFile(c.User.ID)
	} else {
		dir, err = models.GetFileByID(form.ParentID, c.User.ID)
	}

	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusBadRequest, ecode.FileParentNotDirError, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	if !dir.IsDir() {
		c.Error(http.StatusBadRequest, ecode.FileParentNotDirError, err)
		return
	}

	ndir, err := models.CreateDirectory(dir, form.DirectoryName)
	if err != nil {
		if models.IsErrFileAlreadyExist(err) {
			c.Error(http.StatusBadRequest, ecode.FileAlreadyExists, err)
		} else if models.IsErrFileParentNotDirectory(err) {
			c.Error(http.StatusBadRequest, ecode.FileParentNotDirError, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	c.OK(convert.ToFile(ndir))
}

type GetFileInfoForm struct {
	FileID uint `json:"file_id" uri:"file_id" form:"file_id" binding:"required"`
}

func GetFileInfo(c *context.APIContext) {

	fileIDParam := c.Param("file_id")
	fileID, err := strconv.ParseUint(fileIDParam, 10, 64)
	if err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	file, err := models.GetFileByID(uint(fileID), 0)
	if err != nil {
		if models.IsErrFileNotExist(err) {
			c.Error(http.StatusNotFound, ecode.FileNotExist, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	c.OK(convert.ToFile(file))
}
