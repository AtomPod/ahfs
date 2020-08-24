package structs

import "time"

type File struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FileName  string    `json:"filename"`
	FileDir   string    `json:"path"`
	IsDir     bool      `json:"is_dir"`
	FileSize  int64     `json:"size"`
	Owner     uint      `json:"owner"`
	ParentID  uint      `json:"parent_id"`
}
