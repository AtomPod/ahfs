package admin

import (
	"strconv"
	"strings"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/context"
	"github.com/czhj/ahfs/modules/convert"
	api "github.com/czhj/ahfs/modules/structs"
	"github.com/czhj/ahfs/routers/api/v1/utils"
)

func ListsFile(c *context.APIContext) {

	listOptions := utils.GetListOptions(c)
	owner, _ := strconv.Atoi(c.Query("owner"))
	fid, _ := strconv.Atoi(c.Query("fid"))
	typ, _ := strconv.Atoi(c.Query("type"))

	opts := &models.SearchFileOptions{
		ListOptions: listOptions,
		Keyword:     strings.Trim(c.Query("q"), " "),
		Owner:       uint(owner),
		FID:         uint(fid),
		Type:        models.FileType(typ),
	}

	files, maxResult, err := models.SearchFile(opts)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	result := make([]*api.File, len(files))
	for i := range files {
		result[i] = convert.ToFile(files[i])
	}

	c.Header("X-Total-Count", strconv.FormatInt(maxResult, 10))
	c.Header("Access-Control-Expose-Headers", "X-Total-Count")
	c.OK(result)

}
