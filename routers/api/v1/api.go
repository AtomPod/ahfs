package v1

import (
	"net/http"

	"github.com/czhj/ahfs/modules/context"
	ecode "github.com/czhj/ahfs/routers/api/v1/errcode"
	"github.com/czhj/ahfs/routers/api/v1/file"
	"github.com/czhj/ahfs/routers/api/v1/user"
	"github.com/gin-gonic/gin"
)

func requestSignIn() context.APIHandlerFunc {
	return func(c *context.APIContext) {
		if !c.IsSigned {
			c.Error(http.StatusUnauthorized, ecode.UnauthorizedError, "unauthorized error")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RegisterRoutes(e *gin.RouterGroup) {
	v1 := e.Group("/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("/list", context.APIContextWrapper(user.Search))
			users.POST("/", context.APIContextWrapper(user.SignUpPost))
			users.POST("/token", context.APIContextWrapper(user.SignInPost))
			users.PUT("/password", context.APIContextWrapper(user.ResetPasswordPost))
			users.POST("/email_active_code", context.APIContextWrapper(user.RequestActiveEmail))
			users.POST("/reset_password_code", context.APIContextWrapper(user.RequestResetPwdCode))
		}
		userApi := v1.Group("/user")
		{
			userApi.Use(context.APIContextWrapper(requestSignIn()))
			userApi.GET("", context.APIContextWrapper(user.GetAuthenticatedUser))
			userApi.GET("/directory/root", context.APIContextWrapper(file.GetUserRootDirectory))
			userApi.PUT("/info", context.APIContextWrapper(user.ModifyUserInformation))
		}

		files := v1.Group("/files")
		{
			files.Use(context.APIContextWrapper(requestSignIn()))
			files.POST("", context.APIContextWrapper(file.UploadFile))
			files.GET("/:file_id", context.APIContextWrapper(file.DownloadFile))
			files.GET("/:file_id/info", context.APIContextWrapper(file.GetFileInfo))
			files.PUT("/:file_id/name", context.APIContextWrapper(file.RenameFile))
			files.PUT("/:file_id/directory", context.APIContextWrapper(file.MoveFile))
			files.DELETE("/:file_id", context.APIContextWrapper(file.DeleteFile))
		}

		directory := v1.Group("/directory")
		{
			directory.Use(context.APIContextWrapper(requestSignIn()))
			directory.POST("", context.APIContextWrapper(file.CreateDirectory))
			directory.GET("/:file_id", context.APIContextWrapper(file.ReadDirectory))
		}
	}
}
