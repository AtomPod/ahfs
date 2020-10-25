package context

import (
	"fmt"
	"net/http"
	"os"

	"github.com/czhj/ahfs/modules/log"
	"github.com/czhj/ahfs/modules/storage"
	"github.com/czhj/ahfs/routers/api/v1/errcode"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	APIContextKey = "ahfs-api-context-key"
)

type APIHandlerFunc func(*APIContext)

type APIContext struct {
	*Context
}

type APIResult struct {
	Code  errcode.ErrorCode `json:"code"`
	Data  interface{}       `json:"data,omitempty"`
	Error *APIError         `json:"error,omitempty"`
}

type APIError struct {
	URL     string `json:"url"`
	Message string `json:"message"`
}

func (ctx *APIContext) File(filename string, localPath string) {
	file, err := os.Open(localPath)
	if err != nil {
		ctx.InternalServerError(err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error("Failed to close file", zap.String("filepath", localPath), zap.Error(err))
		}
	}()

	fi, err := file.Stat()
	if err != nil {
		ctx.InternalServerError(err)
		return
	}

	ctx.DataFromReader(http.StatusOK, fi.Size(), "application/octet-stream", file,
		map[string]string{
			"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s\"", filename),
		})
}

func (ctx *APIContext) Storage(filename string, id storage.ID, fileStorage storage.Storage) {
	fileObj, err := fileStorage.Read(id)
	if err != nil {
		ctx.InternalServerError(err)
		return
	}
	file := fileObj.Reader

	defer func() {
		if err := file.Close(); err != nil {
			log.Error("Failed to close file", zap.String("id", string(id)), zap.Error(err))
		}
	}()

	ctx.DataFromReader(http.StatusOK, fileObj.Size, "application/octet-stream", file,
		map[string]string{
			"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s\"", filename),
		})
}

func (ctx *APIContext) JSON(status int, code errcode.ErrorCode, obj interface{}) {
	ctx.Context.JSON(status, APIResult{
		Code: code,
		Data: obj,
	})
}

func (ctx *APIContext) OK(data interface{}) {
	ctx.JSON(http.StatusOK, errcode.OK, data)
}

func (ctx *APIContext) Done(status int, data interface{}) {
	ctx.JSON(status, errcode.OK, data)
}

func (ctx *APIContext) NotFound(code errcode.ErrorCode, err error) {
	ctx.Error(http.StatusNotFound, code, err)
}

func (ctx *APIContext) Error(status int, code errcode.ErrorCode, obj interface{}) {

	var message string
	if err, ok := obj.(error); ok {
		message = err.Error()
	} else if obj != nil {
		message = fmt.Sprintf("%s", obj)
	}

	if status == http.StatusInternalServerError {
		log.ErrorWithSkip(1, message, zap.String("status", "InternalServerError"))

		if gin.Mode() == gin.ReleaseMode {
			message = ""
		}
	}

	ctx.Context.JSON(status, APIResult{
		Code: code,
		Error: &APIError{
			Message: message,
		},
	})

}

func (ctx *APIContext) InternalServerError(err error) {
	ctx.Error(http.StatusInternalServerError, errcode.InternalServerError, err)
}

func APIContexter() gin.HandlerFunc {
	return ContextWrapper(func(c *Context) {
		apictx := &APIContext{
			Context: c,
		}
		c.Set(APIContextKey, apictx)
		c.Next()
	})
}

func APIContextWrapper(f func(c *APIContext)) gin.HandlerFunc {
	return ContextWrapper(func(c *Context) {
		value, _ := c.Get(APIContextKey)
		ctx, _ := value.(*APIContext)
		f(ctx)
	})
}
