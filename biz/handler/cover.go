package handler

import (
	"context"

	"simplified-tik-tok/biz/common"

	"github.com/cloudwego/hertz/pkg/app"
)

func Cover(_ context.Context, c *app.RequestContext) {
	name := c.Param("name")
	c.File(common.CoverPath + name)
}
