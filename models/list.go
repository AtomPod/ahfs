package models

import (
	"github.com/czhj/ahfs/modules/setting"
	"github.com/jinzhu/gorm"
)

type ListOptions struct {
	Page     int
	PageSize int
}

func (opts ListOptions) SetEnginePagination(e *gorm.DB) *gorm.DB {
	return e.Limit(opts.PageSize).Offset((opts.Page - 1) * opts.PageSize)
}

func (opts ListOptions) SetDefaultOptions() {
	if opts.PageSize <= 0 {
		opts.PageSize = setting.API.DefaultPagingSize
	}

	if opts.PageSize > setting.API.MaxPagingSize {
		opts.PageSize = setting.API.MaxPagingSize
	}

	if opts.Page <= 0 {
		opts.Page = 1
	}
}

type SearchOrderBy string

func (s SearchOrderBy) String() string {
	return string(s)
}

const (
	SearchOrderByID                      = "id ASC"
	SearchOrderByIDReverse               = "id DESC"
	SearchOrderByOldest                  = "created_at ASC"
	SearchOrderByNewest                  = "created_at DESC"
	SearchOrderByLeastUpdated            = "updated_at ASC"
	SearchOrderByRecentUpdated           = "updated_at DESC"
	SearchOrderByAlphabetically          = "nickname ASC"
	SearchOrderByAlphabeticallyReverse   = "nickname DESC"
	SearchOrderByMaxFileCapacity         = "max_file_capacity ASC"
	SearchOrderByMaxFileCapacityReverse  = "max_file_capacity DESC"
	SearchOrderByUsedFileCapacity        = "used_file_capacity ASC"
	SearchOrderByUsedFileCapacityReverse = "used_file_capacity DESC"
)
