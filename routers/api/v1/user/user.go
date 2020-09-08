package user

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/auth"
	"github.com/czhj/ahfs/modules/cache"
	"github.com/czhj/ahfs/modules/code"

	"github.com/czhj/ahfs/modules/context"
	"github.com/czhj/ahfs/modules/convert"
	"github.com/czhj/ahfs/modules/setting"
	api "github.com/czhj/ahfs/modules/structs"
	"github.com/czhj/ahfs/services/mailer"

	ecode "github.com/czhj/ahfs/routers/api/v1/errcode"
	"github.com/czhj/ahfs/routers/api/v1/utils"
)

func SignInPost(c *context.APIContext) {

	form := &auth.SignInForm{}
	if err := c.Bind(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.IncorrectUserNameOrPwd, err)
		return
	}

	user, err := models.UserSignIn(form.Username, form.Password)
	if err != nil {
		if models.IsErrUserNotExist(err) {
			c.Error(http.StatusNotFound, ecode.IncorrectUserNameOrPwd, err)
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

func RequestActiveEmail(c *context.APIContext) {
	form := &auth.RequestActiveEmailForm{}
	if err := c.Bind(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.EmailFormatError, err)
		return
	}

	exist, err := models.IsEmailUsed(form.Email)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	if exist {
		c.Error(http.StatusBadRequest, ecode.EmailAlreadyExists, fmt.Errorf("email [%s] already exists", form.Email))
		return
	}

	activeCode, err := code.CreateEmailActiveCode(form.Email)
	if err != nil {
		c.InternalServerError(err)
		return
	}
	mailer.SendActiveCodeMail(form.Email, activeCode)
	c.OK(nil)
}

func SignUpPost(c *context.APIContext) {

	form := &auth.SignUpForm{}
	if err := c.Bind(form); err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	ok, err := code.VerifyEmailActiveCode(form.Email, form.EmailVerifyCode)
	if err != nil || !ok {
		if !ok || err == cache.ErrNotFound {
			c.Error(http.StatusNotFound, ecode.EmailActiveCodeError, fmt.Errorf("verification code not found"))
		} else {
			c.InternalServerError(err)
		}
		return
	}

	user := &models.User{
		Username:        form.Username,
		Nickname:        form.Nickname,
		Email:           form.Email,
		Type:            models.UserTypeUser,
		MaxFileCapacity: setting.Service.MaxFileCapacitySize,
		LastLoginAt:     time.Now(),
		Password:        form.Password,
	}

	if err := models.CreateUser(user); err != nil {
		if models.IsErrEmailAlreadyUsed(err) {
			c.Error(http.StatusBadRequest, ecode.EmailAlreadyExists, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	code.RemoveEmailActiveCode(form.Email, form.EmailVerifyCode)
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
	ctx.OK(convert.ToUser(ctx.User, ctx.IsSigned, ctx.User != nil))
}

func GetUserInformation(c *context.APIContext) {
	username := c.Param("username")

	user, err := models.GetUserByUsername(username)
	if err != nil {
		if models.IsErrUserNotExist(err) {
			c.Error(http.StatusNotFound, ecode.UserNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	c.OK(convert.ToUser(user, c.IsSigned, c.User != nil && (c.User.IsAdmin || c.User.ID == user.ID)))
}

func GetUserAvatar(c *context.APIContext) {
	username := c.Param("username")

	user, err := models.GetUserByUsername(username)
	if err != nil {
		if models.IsErrUserNotExist(err) {
			c.Error(http.StatusNotFound, ecode.UserNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	c.Redirect(http.StatusFound, user.AvatarLink())
}
