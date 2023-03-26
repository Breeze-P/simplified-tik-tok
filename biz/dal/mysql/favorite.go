package mysql

import (
	"fmt"

	"simplified-tik-tok/biz/model"
)

func CreateFavorite(favorite *model.Favorite) error {
	return DB.Create(favorite).Error
}

func DeleteFavorite(favorite *model.Favorite) error {
	return DB.Where("user_id = ? AND video_id = ?", favorite.UserID, favorite.VideoID).Delete(&model.Favorite{}).Error
}

func FindVideoIDsByFavoriteActorID(favoriteActorID int64) ([]int64, error) {
	var favoriteList []*model.Favorite
	err := DB.Where("user_id = ?", favoriteActorID).Find(&favoriteList).Error
	if err != nil {
		return nil, err
	}
	var videoID []int64
	fmt.Printf("favorite_list: %v\n", favoriteList)
	for _, fav := range favoriteList {
		videoID = append(videoID, fav.VideoID)
	}
	fmt.Printf("videoId: %v\n", videoID)
	return videoID, nil
}

// 判断是否喜欢
func CheckIsFavorite(userID int64, videoID int64) (bool, error) {
	var favoriteList []*model.Favorite
	err := DB.Where("user_id = ? AND video_id = ?", userID, videoID).Find(&favoriteList).Error
	if err != nil || len(favoriteList) == 0 {
		return false, err
	}
	return true, nil
}
