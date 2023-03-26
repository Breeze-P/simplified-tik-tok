package model

import (
	"gorm.io/gorm"
)

type Relation struct {
	gorm.Model
	FollowerID int64 `json:"follower_id" Column:"follower_id"`
	FollowedID int64 `json:"followed_id" Column:"followed_id"`
}

func (*Relation) TableName() string {
	return "relations"
}
