package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/czhj/ahfs/modules/context"
	"github.com/czhj/ahfs/modules/limiter"
	"github.com/czhj/ahfs/routers/api/v1/admin"
	ecode "github.com/czhj/ahfs/routers/api/v1/errcode"
	"github.com/czhj/ahfs/routers/api/v1/file"
	"github.com/czhj/ahfs/routers/api/v1/user"

	"github.com/gin-contrib/cors"
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

func requestAdmin() context.APIHandlerFunc {
	return func(c *context.APIContext) {
		if c.IsAdmin() {
			c.Error(http.StatusUnauthorized, ecode.UnauthorizedError, "unauthorized error")
			c.Abort()
			return
		}
		c.Next()
	}
}

func requestLimiter() context.APIHandlerFunc {
	return func(c *context.APIContext) {
		if !c.IsSigned {
			addr := c.Request.RemoteAddr
			count, err := limiter.Request(addr, 1, 60*time.Second)
			if err != nil {
				c.InternalServerError(err)
				c.Abort()
				return
			}

			if count > 60 {
				c.Error(http.StatusForbidden, ecode.VisitTooFrequently, fmt.Errorf("Sorry, you are visiting our service too frequent, please try again later."))
				c.Abort()
				return
			}
		}
	}
}

func RegisterRoutes(e *gin.RouterGroup) {
	v1 := e.Group("/v1")
	{
		v1.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"POST", "GET", "PUT", "DELETE"},
			AllowHeaders:    []string{"Origin"},
		}))
		users := v1.Group("/users")
		{
			//users.GET("/list", context.APIContextWrapper(user.Search))
			users.POST("/", context.APIContextWrapper(user.SignUpPost))
			users.POST("/token", context.APIContextWrapper(user.SignInPost))
			users.PUT("/password", context.APIContextWrapper(user.ResetPasswordPost))
			users.POST("/email_active_code", context.APIContextWrapper(user.RequestActiveEmail))
			users.POST("/reset_password_code", context.APIContextWrapper(user.RequestResetPwdCode))
		}

		currentUser := v1.Group("/current_user")
		{
			currentUser.Use(context.APIContextWrapper(requestSignIn()))
			currentUser.GET("", context.APIContextWrapper(user.GetAuthenticatedUser))
			currentUser.GET("/directory/root", context.APIContextWrapper(file.GetUserRootDirectory))
			currentUser.PATCH("/", context.APIContextWrapper(user.EditUser))
			currentUser.PUT("/password", context.APIContextWrapper(user.EditUserPassword))
			currentUser.PUT("/avatar", context.APIContextWrapper(user.UpdateAvatar))
		}

		userGroup := v1.Group("/user")
		{
			userGroup.Use(context.APIContextWrapper(requestLimiter()))
			userGroup.GET("/:username/info", context.APIContextWrapper(user.GetUserInformation))
			userGroup.GET("/:username/avatar", context.APIContextWrapper(user.GetUserAvatar))
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

		adminGroup := v1.Group("/admin", context.APIContextWrapper(requestAdmin()))
		{
			usersAdmin := adminGroup.Group("/users")
			usersAdmin.GET("", context.APIContextWrapper(admin.ListsUser))
			usersAdmin.DELETE("/:username", context.APIContextWrapper(admin.DeleteUser))
			usersAdmin.PATCH("/:username", context.APIContextWrapper(admin.EditUser))

			filesAdmin := adminGroup.Group("/files")
			filesAdmin.GET("", context.APIContextWrapper(admin.ListsFile))
			filesAdmin.GET("/:file_id", context.APIContextWrapper(file.DownloadFile))
			filesAdmin.GET("/:file_id/info", context.APIContextWrapper(file.GetFileInfo))
			filesAdmin.PUT("/:file_id/name", context.APIContextWrapper(file.RenameFile))
			filesAdmin.PUT("/:file_id/directory", context.APIContextWrapper(file.MoveFile))
			filesAdmin.DELETE("/:file_id", context.APIContextWrapper(file.DeleteFile))
		}
	}
}
