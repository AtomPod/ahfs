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
		}
		userApi := v1.Group("/user")
		{
			userApi.GET("", context.APIContextWrapper(user.GetAuthenticatedUser))
			userApi.GET("/:user_id/directory/:directory_id", context.APIContextWrapper(file.ReadDirectory))
		}

		files := v1.Group("/files")
		{
			files.GET("/:file_id/directory", context.APIContextWrapper(file.ReadDirectory))
			files.GET("/:file_id", context.APIContextWrapper(file.DownloadFile))
			files.GET("", context.APIContextWrapper(file.GetUserRootDirectory))

			files.POST("", context.APIContextWrapper(file.UploadFile))
			files.POST("/:file_id/rename", context.APIContextWrapper(file.RenameFile))
			files.POST("/:file_id/move", context.APIContextWrapper(file.MoveFile))

			files.DELETE("/:file_id", context.APIContextWrapper(file.DeleteFile))
		}
	}
}
