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

func CommentList(_ context.Context, c *app.RequestContext) {
	var commentListRequest struct {
		Token   string `json:"token" form:"token" query:"token"`          // 用户鉴权token
		VideoID int64  `json:"video_id" form:"video_id" query:"video_id"` // 视频id
	}

	myUserID := uint(0)
	var err error

	if err = c.BindAndValidate(&commentListRequest); err != nil {
		c.JSON(http.StatusOK, common.NewPrivilegeError("no privilege"))
		return
	}

	if len(commentListRequest.Token) > 0 {
		myUserID, err = mw.GetIDFromTokenString(commentListRequest.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewPrivilegeError("no privilege"))
			return
		}
	}

	commentList, err := mysql.FindCommentsByVideoID(commentListRequest.VideoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewDatabaseError("videoID query wrong"))
		return
	}

	var resCommentList []resbody.Comment

	for _, comment := range commentList {
		videoAuthors, err := mysql.FindUserByID(uint(comment.UserID))
		if err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("commenter info Wrong"))
			return
		}
		videoAuthor := videoAuthors[0]

		isFollow := false
		if myUserID > 0 {
			isFollow, _ = mysql.IsFollow(int64(myUserID), int64(videoAuthor.ID))
		}

		resCommentList = append(resCommentList, resbody.Comment{
			ID: int64(comment.ID),
			User: resbody.User{
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
			Content:    comment.CommentText,
			CreateDate: comment.CreateDate,
		})
	}
	c.JSON(http.StatusOK, utils.H{
		"status_code":  0,
		"status_msg":   "success",
		"comment_list": resCommentList,
	})
}
