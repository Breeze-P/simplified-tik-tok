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

func RelationFollowList(_ context.Context, c *app.RequestContext) {
	var relationFollowListRequest struct {
		UserID int64  `json:"user_id" form:"user_id" query:"user_id"` // 用户id
		Token  string `json:"token" form:"token" query:"token"`       // 用户鉴权token
	}
	if err := c.BindAndValidate(&relationFollowListRequest); err != nil {
		c.JSON(http.StatusOK, common.NewParameterError("feed param form wrong"))
		return
	}

	UserID, err := mw.GetIDFromTokenString(relationFollowListRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, common.NewPrivilegeError("no privilege"))
		return
	}

	followIDs, err := mysql.GetFollowIDsByUserID(int64(UserID))
	if err != nil {
		c.JSON(http.StatusOK, common.NewDatabaseError("follow IDs query wrong"))
		return
	}

	var resUserList []resbody.User

	for _, followID := range followIDs {
		follows, err := mysql.FindUserByID(uint(followID))
		if err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("follow infomation wrong"))
			return
		}
		follow := follows[0]

		resUserList = append(resUserList, resbody.User{
			ID:              int64(follow.ID),
			Name:            follow.Name,
			FollowCount:     follow.FavoriteCount,
			FollowerCount:   follow.FollowCount,
			IsFollow:        true,
			Avatar:          follow.Avatar,
			BackgroundImage: follow.BackgroundImage,
			Signature:       follow.Signature,
			TotalFavorited:  follow.TotalFavorited,
			WorkCount:       follow.WorkCount,
			FavoriteCount:   follow.FavoriteCount,
		})

	}
	c.JSON(http.StatusOK, utils.H{
		"status_code": 0,
		"status_msg":  "success",
		"user_list":   resUserList,
	})
}
