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

func UserInfo(_ context.Context, c *app.RequestContext) {
	var userInfoStruct struct {
		UserID int64  `form:"user_id" json:"user_id" query:"user_id" `
		Token  string `form:"token" json:"token" query:"token"`
	}
	if err := c.BindAndValidate(&userInfoStruct); err != nil {
		c.JSON(http.StatusOK, common.NewParameterError("feed param form wrong"))
		return
	}

	followerID, err := mw.GetIDFromTokenString(userInfoStruct.Token)
	if err != nil {
		c.JSON(http.StatusOK, common.NewPrivilegeError("no privilege"))
		return
	}

	users, err := mysql.FindUserByID(uint(userInfoStruct.UserID)) // 获取作者信息
	if err != nil {
		c.JSON(http.StatusOK, common.NewDatabaseError("userInfo get wrong"))
		return
	}

	user := users[0]

	isFollow, err := mysql.IsFollow(int64(followerID), userInfoStruct.UserID)
	if err != nil {
		c.JSON(http.StatusOK, common.NewDatabaseError("following relation get wrong"))
		return
	}

	userResp := resbody.User{
		ID:              int64(user.ID),
		Name:            user.Name,
		FollowCount:     user.FavoriteCount,
		FollowerCount:   user.FollowCount,
		IsFollow:        isFollow, // Done
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
	c.JSON(http.StatusOK, utils.H{
		"status_code": 0,
		"status_msg":  "success",
		"user":        userResp,
	})
}
