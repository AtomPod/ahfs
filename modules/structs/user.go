package structs

import "time"

type User struct {
	ID              uint      `json:"id"`
	Username        string    `json:"username"`
	NickName        string    `json:"nickname"`
	AvatarURL       string    `json:"avatar_url"`
	Email           string    `json:"email"`
	IsAdmin         bool      `json:"is_admin"`
	CreatedAt       time.Time `json:"created_at"`
	LastLoginAt     time.Time `json:"last_login_at"`
	UsedFileCapcity int64     `json:"used_file_capacity"`
	MaxFileCapacity int64     `json:"max_file_capacity"`
	RootFileID      uint      `json:"root_file_id"`
}
