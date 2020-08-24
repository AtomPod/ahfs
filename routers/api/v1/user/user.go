package user

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/auth"
	"github.com/czhj/ahfs/modules/context"
	"github.com/czhj/ahfs/modules/convert"
	"github.com/czhj/ahfs/modules/setting"
	api "github.com/czhj/ahfs/modules/structs"

	"github.com/czhj/ahfs/routers/api/v1/utils"
)

func SignInPost(c *context.APIContext) {

	form := &auth.SignInForm{}
	if err := c.Bind(form); err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	user, err := models.UserSignIn(form.Username, form.Password)
	if err != nil {
		if models.IsErrUserNotExist(err) {
			c.Error(http.StatusNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	authToken := &models.AuthToken{UserID: user.ID}
	if err := models.CreateAuthToken(authToken); err != nil {
		c.InternalServerError(fmt.Errorf("CreateAuthToken: Unable to create auth token: %v", err))
		return
	}

	result := convert.ToToken(authToken)
	c.OK(result)

}

func SignUpPost(c *context.APIContext) {

	form := &auth.SignUpForm{}
	if err := c.Bind(form); err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	user := &models.User{
		Nickname:        form.Nickname,
		Email:           form.Email,
		Type:            models.UserTypeUser,
		MaxFileCapacity: setting.Service.MaxFileCapacitySize,
		LastLoginAt:     time.Now(),
		Password:        form.Password,
	}

	if err := models.CreateUser(user); err != nil {

		if models.IsErrEmailAlreadyUsed(err) {
			c.Error(http.StatusBadRequest, err)
		} else {
			c.InternalServerError(err)
		}

		return

	}

	authToken := &models.AuthToken{UserID: user.ID}
	if err := models.CreateAuthToken(authToken); err != nil {
		c.InternalServerError(fmt.Errorf("CreateAuthToken: Unable to create auth token: %v", err))
		return
	}

	result := convert.ToToken(authToken)
	c.OK(result)
}

func Search(ctx *context.APIContext) {

	listOptions := utils.GetListOptions(ctx)
	uid, _ := strconv.Atoi(ctx.Query("uid"))

	opts := &models.SearchUserOptions{
		ListOptions: listOptions,
		Keyword:     strings.Trim(ctx.Query("q"), " "),
		UID:         uint(uid),
		Type:        models.UserTypeUser,
	}

	users, maxResult, err := models.SearchUser(opts)
	if err != nil {
		ctx.InternalServerError(err)
		return
	}

	result := make([]*api.User, len(users))
	for i := range users {
		result[i] = convert.ToUser(users[i], ctx.IsSigned, ctx.User != nil && ctx.User.IsAdmin)
	}

	ctx.Header("X-Total-Count", strconv.FormatInt(maxResult, 10))
	ctx.Header("Access-Control-Expose-Headers", "X-Total-Count")
	ctx.OK(result)
}

func GetAuthenticatedUser(ctx *context.APIContext) {
	if ctx.User == nil {
		ctx.Error(http.StatusUnauthorized, nil)
		return
	}
	ctx.OK(convert.ToUser(ctx.User, ctx.IsSigned, ctx.User != nil))
}
