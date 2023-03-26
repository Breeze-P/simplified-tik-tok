package model

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	UserID  int64 `json:"user_id" Column:"user_id"`
	VideoID int64 `json:"video_id" Column:"video_id"`
}

func (*Favorite) TableName() string {
	return "favorites"
}
