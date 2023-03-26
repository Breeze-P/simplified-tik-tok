package model

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Createtime    int64  `json:"create_time" Column:"createtime"`
	AuthorID      uint   `json:"author_id" Column:"author_id"`
	Title         string `json:"title" Column:"title"`
	Filepath      string `json:"file_path" Column:"filepath"`
	Playurl       string `json:"play_url" Column:"playurl"`
	Coverurl      string `json:"cover_url" Column:"coverurl"`
	Favoritecount int64  `json:"favorite_count" Column:"favoritecount"`
	Commentcount  int64  `json:"comment_count" Column:"commentcount"`
}

func (v *Video) TableName() string {
	return "videos"
}
