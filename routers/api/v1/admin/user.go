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
