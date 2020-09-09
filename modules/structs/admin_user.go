package structs

type EditUserOption struct {
	Nickname           *string `form:"nickname" json:"nickname" binding:"omitempty,gt=1,lt=64"`
	Admin              *bool   `form:"is_admin" json:"is_admin" binding:"omitempty"`
	Active             *bool   `form:"is_active" json:"is_active" binding:"omitempty"`
	MustChangePassword *bool   `form:"must_change_password" json:"must_change_password" binding:"omitempty"`
	MaxFileCapacity    *int64  `form:"max_file_capacity" json:"max_file_capacity" binding:"omitempty"`
}
