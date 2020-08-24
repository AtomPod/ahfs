package convert

import (
	"time"

	"github.com/czhj/ahfs/models"
	api "github.com/czhj/ahfs/modules/structs"
)

func ToUser(u *models.User, signed bool, authed bool) *api.User {
	user := &api.User{
		Username:  u.Username,
		NickName:  u.Nickname,
		AvatarURL: u.AvatarLink(),
		CreatedAt: u.CreatedAt,
	}

	if signed && authed {
		user.Email = u.Email
	}

	if authed {
		user.ID = u.ID
		user.IsAdmin = u.IsAdmin
		user.LastLoginAt = u.LastLoginAt
		user.MaxFileCapacity = u.MaxFileCapacity
		user.UsedFileCapcity = u.UsedFileCapacity
	}
	return user
}

func ToToken(t *models.AuthToken) *api.Token {
	return &api.Token{
		Token:     t.Code,
		ExpiresIn: uint64(time.Duration(7 * 24 * time.Hour).Seconds()),
	}
}

func ToFile(f *models.File) *api.File {
	return &api.File{
		ID:        f.ID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		FileName:  f.FileName,
		IsDir:     f.IsDir(),
		Owner:     f.Owner,
		FileSize:  f.FileSize,
		ParentID:  f.ParentID,
	}
}
