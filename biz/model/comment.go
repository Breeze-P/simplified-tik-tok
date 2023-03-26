package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID      int64  `json:"user_id" Column:"user_id"`
	VideoID     int64  `json:"video_id" Column:"video_id"`
	CommentText string `json:"conmment_text" Column:"comment_text"`
	CreateDate  string `json:"create_date" Column:"create_date"`
}

func (*Comment) TableName() string {
	return "comment"
}
