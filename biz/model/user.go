package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" Column:"username"`
	Password string `json:"password" Column:"password"`

	Avatar          string `json:"avatar" Column:"avatar"`                     // 用户头像
	BackgroundImage string `json:"background_image" Column:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count" Column:"favorite_count"`     // 喜欢数
	FollowCount     int64  `json:"follow_count" Column:"follow_count"`         // 关注总数
	FollowerCount   int64  `json:"follower_count" Column:"follower_count"`     // 粉丝总数
	Name            string `json:"name" Column:"name"`                         // 用户名称
	Signature       string `json:"signature" Column:"favorite_count"`          // 个人简介
	TotalFavorited  int64  `json:"total_favorited" Column:"total_favorited"`   // 获赞数量
	WorkCount       int64  `json:"work_count" Column:"work_count"`             // 作品数

	Publishedvideos []Video `gorm:"foreignKey:AuthorID;"`
}

func (u *User) TableName() string {
	return "users"
}
