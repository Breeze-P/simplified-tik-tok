package handler

import (
	"context"
	"fmt"

	"simplified-tik-tok/biz/model"
	"simplified-tik-tok/biz/mw"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// Ping .
func Ping(_ context.Context, c *app.RequestContext) {
	user, _ := c.Get(mw.IdentityKey)
	if user == nil {
		c.JSON(200, utils.H{
			"status_code": 200,
			"status_msg":  fmt.Sprintf("pong"),
		})
		return
	}
	c.JSON(200, utils.H{
		"status_code": 200,
		"status_msg":  fmt.Sprintf("username:%v", user.(*model.User).Username),
	})
}
