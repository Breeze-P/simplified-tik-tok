package mysql

import (
	"simplified-tik-tok/biz/model"
)

func CreateComment(comment *model.Comment) error {
	return DB.Create(comment).Error
}

func DeleteCommentByID(id int64) error {
	return DB.Delete(&model.Comment{}, id).Error
}

func FindCommentsByVideoID(id int64) ([]*model.Comment, error) {
	res := make([]*model.Comment, 0)
	err := DB.Where("video_id = ?", id).Find(&res).Error
	return res, err
}
