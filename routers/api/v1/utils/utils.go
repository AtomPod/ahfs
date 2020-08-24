package utils

import (
	"strconv"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/context"
	"github.com/czhj/ahfs/modules/setting"
)

func GetListOptions(ctx *context.APIContext) models.ListOptions {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("limit"))

	if pageSize <= 0 {
		pageSize = setting.API.DefaultPagingSize
	}

	if pageSize > setting.API.MaxPagingSize {
		pageSize = setting.API.MaxPagingSize
	}

	return models.ListOptions{
		Page:     page,
		PageSize: pageSize,
	}
}
