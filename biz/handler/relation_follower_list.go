package handler

import (
	"context"
	"net/http"

	"simplified-tik-tok/biz/common"
	"simplified-tik-tok/biz/dal/mysql"
	"simplified-tik-tok/biz/mw"
	"simplified-tik-tok/biz/resbody"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func RelationFollowerList(_ context.Context, c *app.RequestContext) {
	var relationFollowerListRequest struct {
		UserID int64  `json:"user_id" form:"user_id" query:"user_id"` // 用户id
		Token  string `json:"token" form:"token" query:"token"`       // 用户鉴权token
	}
	if err := c.BindAndValidate(&relationFollowerListRequest); err != nil {
		c.JSON(http.StatusOK, common.NewParameterError("feed param form wrong"))
		return
	}

	UserID, err := mw.GetIDFromTokenString(relationFollowerListRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, common.NewPrivilegeError("no privilege"))
		return
	}

	fanIDs, err := mysql.GetFanIDsByUserID(int64(UserID))
	if err != nil {
		c.JSON(http.StatusOK, common.NewDatabaseError("Fan IDs query wrong"))
		return
	}

	var resUserList []resbody.User

	for _, fanID := range fanIDs {
		fans, err := mysql.FindUserByID(uint(fanID))
		if err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("fan infomation wrong"))
			return
		}
		fan := fans[0]

		isFollow, _ := mysql.IsFollow(int64(UserID), int64(fan.ID))

		resUserList = append(resUserList, resbody.User{
			ID:              int64(fan.ID),
			Name:            fan.Name,
			FollowCount:     fan.FavoriteCount,
			FollowerCount:   fan.FollowCount,
			IsFollow:        isFollow,
			Avatar:          fan.Avatar,
			BackgroundImage: fan.BackgroundImage,
			Signature:       fan.Signature,
			TotalFavorited:  fan.TotalFavorited,
			WorkCount:       fan.WorkCount,
			FavoriteCount:   fan.FavoriteCount,
		})

	}
	c.JSON(http.StatusOK, utils.H{
		"status_code": 0,
		"status_msg":  "success",
		"user_list":   resUserList,
	})
}
