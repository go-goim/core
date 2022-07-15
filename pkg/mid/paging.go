package mid

import (
	"github.com/gin-gonic/gin"

	"github.com/go-goim/core/pkg/web"
)

const (
	pagingKey = "paging"
)

func PagingHandler(c *gin.Context) {
	req := &web.Paging{}
	_ = c.ShouldBindQuery(req)
	// set default page size
	if req.Page == 0 {
		req.Page = 1
	}

	if req.PageSize == 0 {
		req.PageSize = 10
	}

	c.Set(pagingKey, req)
}

func GetPaging(c *gin.Context) *web.Paging {
	v, ok := c.Get(pagingKey)
	if !ok {
		return &web.Paging{}
	}

	paging, ok := v.(*web.Paging)
	if !ok {
		return &web.Paging{}
	}

	return paging
}
