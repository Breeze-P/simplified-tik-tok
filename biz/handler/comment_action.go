package handler

import (
	"context"
	"net/http"
	"time"

	"simplified-tik-tok/biz/common"
	"simplified-tik-tok/biz/dal/mysql"
	"simplified-tik-tok/biz/model"
	"simplified-tik-tok/biz/mw"
	"simplified-tik-tok/biz/resbody"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func CommentAction(ctx context.Context, c *app.RequestContext) {
	var commentActionRequest struct {
		Token       string `json:"token" form:"token" query:"token"`                      // 用户鉴权token
		VideoID     int64  `json:"video_id" form:"video_id" query:"video_id"`             // 视频id
		ActionType  int32  `json:"action_type" form:"action_type" query:"action_type"`    // 1-发布评论，2-删除评论
		CommentText string `json:"comment_text" form:"comment_text" query:"comment_text"` // 用户删除的评论id，在action_type=2的时候使用
		CommentID   int64  `json:"comment_id" form:"comment_id" query:"comment_id"`
	}
	if err := c.BindAndValidate(&commentActionRequest); err != nil {
		c.JSON(http.StatusOK, common.NewParameterError("param form wrong"))
		return
	}

	UserID, err := mw.GetIDFromTokenString(commentActionRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, common.NewPrivilegeError("no Privilege"))
		return
	}

	switch commentActionRequest.ActionType {
	case 1:
		currentComment := model.Comment{
			UserID:      int64(UserID),
			VideoID:     commentActionRequest.VideoID,
			CommentText: commentActionRequest.CommentText,
			CreateDate:  time.Now().Format("2006-01-02"),
		}
		if err := mysql.CreateComment(&currentComment); err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("create comment wrong"))
			return
		}

		commentUsers, err := mysql.FindUserByID(UserID)
		if err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("comment message return wrong"))
			return
		}
		commentUser := commentUsers[0]

		video, err := mysql.FindVideoByID(commentActionRequest.VideoID)
		if err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("comment message return wrong"))
			return
		}

		isFollow, err := mysql.IsFollow(int64(UserID), int64(video.AuthorID))
		if err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("comment message return wrong"))
			return
		}

		c.JSON(200, utils.H{
			"status_code": 0,
			"status_msg":  "success",
			"comment": resbody.Comment{
				ID: int64(currentComment.ID),
				User: resbody.User{
					ID:              int64(commentUser.ID),
					Name:            commentUser.Name,
					FollowCount:     commentUser.FavoriteCount,
					FollowerCount:   commentUser.FollowCount,
					IsFollow:        isFollow,
					Avatar:          commentUser.Avatar,
					BackgroundImage: commentUser.BackgroundImage,
					Signature:       commentUser.Signature,
					TotalFavorited:  commentUser.TotalFavorited,
					WorkCount:       commentUser.WorkCount,
					FavoriteCount:   commentUser.FavoriteCount,
				},
				Content:    currentComment.CommentText,
				CreateDate: currentComment.CreateDate,
			},
		})
	case 2:
		if err := mysql.DeleteCommentByID(commentActionRequest.CommentID); err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("delete comment wrong"))
			return
		}
		if err := mysql.DecVideoCommentCountByID(uint(commentActionRequest.VideoID)); err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("delete comment count wrong"))
			return
		}
		c.JSON(200, utils.H{
			"status_code": 0,
			"status_msg":  "success",
		})
	}
}
