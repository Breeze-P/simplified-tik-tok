package handler

import (
	"context"
	"fmt"
	"net/http"

	"simplified-tik-tok/biz/common"
	"simplified-tik-tok/biz/dal/mysql"
	"simplified-tik-tok/biz/mw"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func Feed(_ context.Context, c *app.RequestContext) {
	var feedStruct struct {
		LastestTime int64  `form:"latest_time" json:"latest_time" query:"latest_time" `
		Number      int    `form:"number" json:"number" query:"number" `
		Token       string `form:"token" json:"token" query:"token"`
	}
	if err := c.BindAndValidate(&feedStruct); err != nil {
		c.JSON(http.StatusOK, common.NewParameterError("feed param form wrong"))
		return
	}

	var myUserID uint
	var err error

	if len(feedStruct.Token) == 0 {
		myUserID = 0
	} else {
		myUserID, err = mw.GetIDFromTokenString(feedStruct.Token)
		fmt.Println(err)
		if err != nil {
			c.JSON(http.StatusOK, common.NewPrivilegeError("no privilege"))
			return
		}
	}

	resVideoList, nextTime, err := mysql.GetFeedList(feedStruct.LastestTime, 10, myUserID)
	if err != nil {
		c.JSON(http.StatusOK, common.NewDatabaseError("feed list query failed"))
		return
	}
	c.JSON(http.StatusOK, utils.H{
		"status_code": 0,
		"status_msg":  "success",
		"video_list":  resVideoList,
		"next_time":   nextTime,
	})
}
