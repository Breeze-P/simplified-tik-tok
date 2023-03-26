package handler

import (
	"context"
	"fmt"
	"net/http"

	"simplified-tik-tok/biz/common"
	"simplified-tik-tok/biz/dal/mysql"
	"simplified-tik-tok/biz/mw"
	"simplified-tik-tok/biz/resbody"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func FavoriteList(_ context.Context, c *app.RequestContext) {
	var favoriteListRequest struct {
		UserID int64  `form:"user_id" json:"user_id" query:"user_id"` // 用户id
		Token  string `form:"token" json:"token" query:"token"`       // 用户鉴权token
	}
	if err := c.BindAndValidate(&favoriteListRequest); err != nil {
		c.JSON(http.StatusOK, common.NewParameterError("wrong param form"))
		return
	}

	fmt.Println("hello?")

	myUserID := uint(0)
	var err error

	if len(favoriteListRequest.Token) != 0 {
		myUserID, err = mw.GetIDFromTokenString(favoriteListRequest.Token)
		if err != nil {
			c.JSON(http.StatusOK, common.NewPrivilegeError("no privilege"))
			return
		}
	}

	fmt.Println(myUserID, favoriteListRequest.UserID)
	videoIDs, err := mysql.FindVideoIDsByFavoriteActorID(favoriteListRequest.UserID) // 获取作者喜欢的视频ID
	if err != nil {
		c.JSON(http.StatusOK, common.NewDatabaseError("videoID not found"))
		return
	}

	var resVideoList []resbody.Video

	for _, videoID := range videoIDs {
		video, err := mysql.FindVideoByID(videoID)
		if err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("videoID not found"))
			return
		}
		videoAuthors, err := mysql.FindUserByID(video.AuthorID)
		if err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("authorID not found"))
			return
		}
		videoAuthor := videoAuthors[0]

		isFollow := false
		if myUserID > 0 {
			isFollow, err = mysql.IsFollow(int64(myUserID), int64(videoAuthor.ID))
			if err != nil {
				c.JSON(http.StatusOK, common.NewDatabaseError("relation query error"))
				return
			}
		}

		isFavorite := false
		if myUserID > 0 {
			isFavorite, err = mysql.IsFollow(int64(myUserID), int64(videoAuthor.ID))
			if err != nil {
				c.JSON(http.StatusOK, common.NewDatabaseError("relation query error"))
				return
			}
		}
		resVideoList = append(resVideoList, resbody.Video{
			ID: int64(video.ID),
			Author: resbody.User{
				ID:              int64(videoAuthor.ID),
				Name:            videoAuthor.Name,
				FollowCount:     videoAuthor.FavoriteCount,
				FollowerCount:   videoAuthor.FollowCount,
				IsFollow:        isFollow,
				Avatar:          videoAuthor.Avatar,
				BackgroundImage: videoAuthor.BackgroundImage,
				Signature:       videoAuthor.Signature,
				TotalFavorited:  videoAuthor.TotalFavorited,
				WorkCount:       videoAuthor.WorkCount,
				FavoriteCount:   videoAuthor.FavoriteCount,
			},
			PlayURL:       video.Playurl,
			CoverURL:      video.Coverurl,
			FavoriteCount: video.Favoritecount,
			CommentCount:  video.Commentcount,
			IsFavorite:    isFavorite,
			Title:         video.Title,
		})
	}
	c.JSON(http.StatusOK, utils.H{
		"status_code": 0,
		"status_msg":  "success",
		"video_list":  resVideoList,
	})
}
