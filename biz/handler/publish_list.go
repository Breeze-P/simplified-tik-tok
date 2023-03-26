package handler

import (
	"context"
	"net/http"

	"simplified-tik-tok/biz/common"
	"simplified-tik-tok/biz/dal/mysql"
	"simplified-tik-tok/biz/mw"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func PublishList(_ context.Context, c *app.RequestContext) {
	var listStruct struct {
		UserID int64  `form:"user_id" json:"user_id" query:"user_id" `
		Token  string `form:"token" json:"token" query:"token"`
	}

	if err := c.BindAndValidate(&listStruct); err != nil {
		c.JSON(http.StatusOK, common.NewParameterError("feed param form wrong"))
		return
	}

	myUserID := uint(0)
	var err error

	if len(listStruct.Token) != 0 {
		myUserID, err = mw.GetIDFromTokenString(listStruct.Token)
		if err != nil {
			c.JSON(http.StatusOK, common.NewPrivilegeError("no privilege"))
			return
		}
	}

	resVideoList, err := mysql.GetPublishList(uint(listStruct.UserID), myUserID)
	if err != nil {
		c.JSON(http.StatusOK, common.NewDatabaseError("video info search wrong"))
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"status_code": 0,
		"status_msg":  "success",
		"video_list":  resVideoList,
	})
}
