package mysql

import (
	"simplified-tik-tok/biz/model"
	"simplified-tik-tok/biz/resbody"

	"gorm.io/gorm"
)

// 创建视频记录
func CreateVideo(video *model.Video) error {
	return DB.Create(video).Error
}

func FindVideoByID(id int64) (*model.Video, error) {
	var res *model.Video
	if err := DB.Where("ID = ?", id).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func FindVideosByAuthorID(authorID uint) ([]*model.Video, error) {
	var videoList []*model.Video
	err := DB.Where("author_id= ?", authorID).Find(&videoList).Error
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

func FindNVideosBeforeTime(timeStamp int64, num int) ([]*model.Video, error) {
	var videoList []*model.Video
	err := DB.Where("createtime< ?", timeStamp).Order("createtime desc").Limit(num).Find(&videoList).Error
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

// 封装视频列表的业务逻辑，isFollow写死了
func GetPublishList(userID uint, myUserID uint) ([]resbody.Video, error) {
	videoList, err := FindVideosByAuthorID(uint(userID)) // 获取所有视频
	if err != nil {
		return nil, err
	}
	videoAuthors, err := FindUserByID(uint(userID)) // 获取作者信息
	if err != nil {
		return nil, err
	}

	videoAuthor := videoAuthors[0]

	isFollow := false
	if myUserID > 0 {
		isFollow, err = IsFollow(int64(myUserID), int64(userID))
		if err != nil {
			return nil, err
		}
	}

	author := resbody.User{
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
	}

	resVideoList := make([]resbody.Video, len(videoList), cap(videoList)) // 创建response所需要的信息
	for i := 0; i < len(videoList); i++ {
		resVideoList[i].Author = author
		resVideoList[i].CommentCount = videoList[i].Commentcount
		resVideoList[i].CoverURL = videoList[i].Coverurl
		resVideoList[i].FavoriteCount = videoList[i].Favoritecount
		resVideoList[i].ID = int64(videoList[i].ID)
		resVideoList[i].PlayURL = videoList[i].Playurl
		resVideoList[i].Title = videoList[i].Title
		isFavorite := false
		if myUserID > 0 {
			isFavorite, err = CheckIsFavorite(int64(myUserID), resVideoList[i].ID)
			if err != nil {
				return nil, err
			}
		}

		resVideoList[i].IsFavorite = isFavorite
	}

	return resVideoList, nil
}

// 联表查询影响效率，封装一下
func GetFeedList(timeStamp int64, num int, myUserID uint) ([]resbody.Video, int64, error) {
	videoList, err := FindNVideosBeforeTime(timeStamp, num)
	if err != nil || len(videoList) == 0 {
		return nil, 0, err
	}

	authorCache := make(map[uint]resbody.User)

	resVideoList := make([]resbody.Video, len(videoList), cap(videoList))
	for idx, video := range videoList {

		author, ok := authorCache[video.AuthorID]
		if !ok {
			res, err := FindUserByID(video.AuthorID)
			if err != nil {
				return nil, 0, err
			}
			isFollow, err := IsFollow(int64(myUserID), int64(res[0].ID))
			if err != nil {
				return nil, 0, err
			}
			author = resbody.User{
				ID:              int64(res[0].ID),
				Name:            res[0].Name,
				FollowCount:     res[0].FavoriteCount,
				FollowerCount:   res[0].FollowCount,
				IsFollow:        isFollow,
				Avatar:          res[0].Avatar,
				BackgroundImage: res[0].BackgroundImage,
				Signature:       res[0].Signature,
				TotalFavorited:  res[0].TotalFavorited,
				WorkCount:       res[0].WorkCount,
				FavoriteCount:   res[0].FavoriteCount,
			}
			authorCache[res[0].ID] = author
		}

		resVideoList[idx].Author = author
		resVideoList[idx].CommentCount = video.Commentcount
		resVideoList[idx].CoverURL = video.Coverurl
		resVideoList[idx].FavoriteCount = video.Favoritecount
		resVideoList[idx].ID = int64(video.ID)
		isFavorite := false
		if myUserID > 0 {
			isFavorite, err = CheckIsFavorite(int64(myUserID), resVideoList[idx].ID)
			if err != nil {
				return nil, 0, err
			}
		}

		resVideoList[idx].IsFavorite = isFavorite
		resVideoList[idx].PlayURL = video.Playurl
		resVideoList[idx].Title = video.Title
	}

	return resVideoList, videoList[len(videoList)-1].Createtime, nil
}

func AddVideoFavoriteCountByID(id uint) error {
	return DB.Model(&model.Video{}).Where("id = ?", id).Update("favoritecount", gorm.Expr("favoritecount + ?", 1)).Error
}

func DecVideoFavoriteCountByID(id uint) error {
	return DB.Model(&model.Video{}).Where("id = ?", id).Update("favoritecount", gorm.Expr("favoritecount - ?", 1)).Error
}

func AddVideoCommentCountByID(id uint) error {
	return DB.Model(&model.Video{}).Where("id = ?", id).Update("commentcount", gorm.Expr("favoritecount + ?", 1)).Error
}

func DecVideoCommentCountByID(id uint) error {
	return DB.Model(&model.Video{}).Where("id = ?", id).Update("commentcount", gorm.Expr("favoritecount - ?", 1)).Error
}
