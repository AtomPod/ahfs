package v1

import (
	"github.com/czhj/ahfs/modules/context"
	"github.com/czhj/ahfs/routers/api/v1/file"
	"github.com/czhj/ahfs/routers/api/v1/user"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(e *gin.RouterGroup) {
	v1 := e.Group("/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("/search", context.APIContextWrapper(user.Search))
			users.POST("/signin", context.APIContextWrapper(user.SignInPost))
			users.POST("/signup", context.APIContextWrapper(user.SignUpPost))
			users.POST("/password", context.APIContextWrapper(user.ResetPasswordPost))
			users.POST("/email_active_code", context.APIContextWrapper(user.RequestActiveEmail))
			users.POST("/reset_password_code", context.APIContextWrapper(user.RequestResetPwdCode))
		}
		userApi := v1.Group("/user")
		{
			userApi.GET("", context.APIContextWrapper(user.GetAuthenticatedUser))
		}

		files := v1.Group("/files")
		{
			files.GET("/:file_id", context.APIContextWrapper(file.DownloadFile))
			files.POST("", context.APIContextWrapper(file.UploadFile))
			files.POST("/:file_id/rename", context.APIContextWrapper(file.RenameFile))
			files.POST("/:file_id/move", context.APIContextWrapper(file.MoveFile))
			files.DELETE("/:file_id", context.APIContextWrapper(file.DeleteFile))
		}

		directory := v1.Group("/directory")
		{
			directory.POST("", context.APIContextWrapper(file.CreateDirectory))
			directory.GET("/:file_id", context.APIContextWrapper(file.ReadDirectory))
		}
	}
}
