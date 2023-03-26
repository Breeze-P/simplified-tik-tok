package mysql

import (
	"strconv"

	"simplified-tik-tok/biz/model"

	"gorm.io/gorm"
)

func CreateUsers(users []*model.User) error {
	return DB.Create(users).Error
}

func FindUserByName(username string) ([]*model.User, error) {
	res := make([]*model.User, 0)
	if err := DB.Where(DB.Or("username = ?", username)).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func FindUserByID(id uint) ([]*model.User, error) {
	res := make([]*model.User, 0)
	if err := DB.Where("ID = ?", id).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func FindUserByIDInStr(id string) ([]*model.User, error) {
	idTemp, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return FindUserByID(uint(idTemp))
}

func CheckUser(username, password string) ([]*model.User, error) {
	res := make([]*model.User, 0)
	if err := DB.Where(DB.Or("username = ?", username)).Where("password = ?", password).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// func FindUserByID(user_id string) (*model.User, error) {
// 	var res *model.User
// 	if err := DB.Where(DB.Or("user_id = ?", user_id)).
// 		Find(&res).Error; err != nil {
// 		return nil, err
// 	}
// 	return res, nil
// }

func AddUserFavoriteCountByID(id uint) error {
	return DB.Model(&model.User{}).Where("id = ?", id).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
}

func DecUserFavoriteCountByID(id uint) error {
	return DB.Model(&model.User{}).Where("id = ?", id).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
}

func AddUserTotalFavoritedByID(id uint) error {
	return DB.Model(&model.User{}).Where("id = ?", id).Update("total_favorited", gorm.Expr("total_favorited + ?", 1)).Error
}

func DecUserTotalFavoritedByID(id uint) error {
	return DB.Model(&model.User{}).Where("id = ?", id).Update("total_favorited", gorm.Expr("total_favorited - ?", 1)).Error
}

func AddUserFollowCountByID(id uint) error {
	return DB.Model(&model.User{}).Where("id = ?", id).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error
}

func DecUserFollowCountByID(id uint) error {
	return DB.Model(&model.User{}).Where("id = ?", id).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error
}

func AddUserFansCountByID(id uint) error {
	return DB.Model(&model.User{}).Where("id = ?", id).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error
}

func DecUserFansCountByID(id uint) error {
	return DB.Model(&model.User{}).Where("id = ?", id).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error
}
