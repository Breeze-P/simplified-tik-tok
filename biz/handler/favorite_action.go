package handler

import (
	"context"
	"net/http"

	"simplified-tik-tok/biz/common"
	"simplified-tik-tok/biz/dal/mysql"
	"simplified-tik-tok/biz/model"
	"simplified-tik-tok/biz/mw"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func FavoriteAction(_ context.Context, c *app.RequestContext) {
	var favoriteActionRequest struct {
		Token      string `form:"token" json:"token" query:"token"`                   // 用户鉴权token
		VideoID    int64  `form:"video_id" json:"video_id" query:"video_id"`          // 视频id
		ActionType int32  `form:"action_type" json:"action_type" query:"action_type"` // 1-点赞，2-取消点赞
	}
	if err := c.BindAndValidate(&favoriteActionRequest); err != nil {
		c.JSON(http.StatusOK, common.NewParameterError("wrong param form"))
		return
	}
	UserID, err := mw.GetIDFromTokenString(favoriteActionRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, common.NewPrivilegeError("no privilege"))
		return
	}
	switch favoriteActionRequest.ActionType {
	case 1:
		if err := mysql.CreateFavorite(&model.Favorite{
			UserID:  int64(UserID),
			VideoID: favoriteActionRequest.VideoID,
		}); err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("create favorite wrong"))
			return
		}

		mysql.AddUserFavoriteCountByID(UserID)
		mysql.AddVideoFavoriteCountByID(uint(favoriteActionRequest.VideoID))
		video, _ := mysql.FindVideoByID(favoriteActionRequest.VideoID)
		mysql.AddUserTotalFavoritedByID(video.AuthorID)
	case 2:
		if err := mysql.DeleteFavorite(&model.Favorite{
			UserID:  int64(UserID),
			VideoID: favoriteActionRequest.VideoID,
		}); err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("delete favorite wrong"))
			return
		}
		mysql.DecUserFavoriteCountByID(UserID)
		mysql.DecVideoFavoriteCountByID(uint(favoriteActionRequest.VideoID))
		video, _ := mysql.FindVideoByID(favoriteActionRequest.VideoID)
		mysql.DecUserTotalFavoritedByID(video.AuthorID)
	default:
		c.JSON(http.StatusOK, common.NewParameterError("action-type wrong"))
		return
	}
	c.JSON(200, utils.H{
		"status_code": 0,
		"status_msg":  "success",
	})
}
