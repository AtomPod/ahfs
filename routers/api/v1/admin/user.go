package admin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/context"
	"github.com/czhj/ahfs/modules/convert"
	api "github.com/czhj/ahfs/modules/structs"
	ecode "github.com/czhj/ahfs/routers/api/v1/errcode"
	"github.com/czhj/ahfs/routers/api/v1/utils"
)

func DeleteUser(c *context.APIContext) {
	username := c.Param("username")
	user, err := models.GetUserByUsername(username)

	if err != nil {
		if models.IsErrUserNotExist(err) {
			c.NotFound(ecode.UsernameNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	if err := models.DeleteUser(user); err != nil {
		if models.IsErrUserNotExist(err) {
			c.NotFound(ecode.UsernameNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	c.Done(http.StatusNotFound, nil)
}

func ListsUser(c *context.APIContext) {

	listOptions := utils.GetListOptions(c)
	uid, _ := strconv.Atoi(c.Query("uid"))

	opts := &models.SearchUserOptions{
		ListOptions: listOptions,
		Keyword:     strings.Trim(c.Query("q"), " "),
		UID:         uint(uid),
		Type:        models.UserTypeUser,
	}

	users, maxResult, err := models.SearchUser(opts)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	result := make([]*api.User, len(users))
	for i := range users {
		result[i] = convert.ToUser(users[i], c.IsSigned, c.User.IsAdmin)
	}

	c.Header("X-Total-Count", strconv.FormatInt(maxResult, 10))
	c.Header("Access-Control-Expose-Headers", "X-Total-Count")
	c.OK(result)

}

func EditUser(c *context.APIContext) {

	username := c.Param("username")

	opts := &api.EditUserOption{}
	if err := c.ShouldBind(opts); err != nil {
		c.Error(http.StatusBadRequest, ecode.ParameterFormatError, err)
		return
	}

	user, err := models.GetUserByUsername(username)
	if err != nil {
		if models.IsErrUserNotExist(err) {
			c.NotFound(ecode.UserNotFound, err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	if opts.Active != nil {
		user.IsActive = *opts.Active
	}

	if opts.Admin != nil {
		user.IsAdmin = *opts.Admin
	}

	if opts.MustChangePassword != nil {
		user.MustChangePassword = *opts.MustChangePassword
	}

	if opts.Nickname != nil {
		user.Nickname = *opts.Nickname
	}

	if opts.MaxFileCapacity != nil {
		user.MaxFileCapacity = *opts.MaxFileCapacity
	}

	if err := models.SaveUser(user); err != nil {
		c.InternalServerError(err)
		return
	}

	c.OK(convert.ToUser(user, c.IsSigned, c.IsAdmin()))
}
