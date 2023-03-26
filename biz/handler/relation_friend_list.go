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

func RelationFriendList(_ context.Context, c *app.RequestContext) {
	var relationFriendListRequest struct {
		UserID int64  `json:"user_id" form:"user_id" query:"user_id"` // 用户id
		Token  string `json:"token" form:"token" query:"token"`       // 用户鉴权token
	}
	if err := c.BindAndValidate(&relationFriendListRequest); err != nil {
		c.JSON(http.StatusOK, common.NewParameterError("feed param form wrong"))
		return
	}

	UserID, err := mw.GetIDFromTokenString(relationFriendListRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, common.NewPrivilegeError("no privilege"))
		return
	}

	friendIDs, err := mysql.GetFriendIDsByUserID(int64(UserID))
	if err != nil {
		c.JSON(http.StatusOK, common.NewDatabaseError("friend IDs query wrong"))
		return
	}

	var resUserList []resbody.FriendUser

	for _, friendID := range friendIDs {
		friends, err := mysql.FindUserByID(uint(friendID))
		if err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("friend infomation wrong"))
			return
		}
		friend := friends[0]

		resUserList = append(resUserList, resbody.FriendUser{
			ID:              int64(friend.ID),
			Name:            friend.Name,
			FollowCount:     friend.FavoriteCount,
			FollowerCount:   friend.FollowCount,
			IsFollow:        true,
			Avatar:          friend.Avatar,
			BackgroundImage: friend.BackgroundImage,
			Signature:       friend.Signature,
			TotalFavorited:  friend.TotalFavorited,
			WorkCount:       friend.WorkCount,
			FavoriteCount:   friend.FavoriteCount,

			Message: "test message",
			MsgType: 0,
		})

	}
	c.JSON(http.StatusOK, utils.H{
		"status_code": 0,
		"status_msg":  "success",
		"user_list":   resUserList,
	})
}
