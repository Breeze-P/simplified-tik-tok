package mysql

import (
	"fmt"

	"simplified-tik-tok/biz/model"
)

// 增加一条关系
func CreateRelation(followerID int64, followedID int64) error {
	if followedID == followerID {
		return nil
	}
	return DB.Create([]*model.Relation{
		{
			FollowerID: followerID,
			FollowedID: followedID,
		},
	}).Error
}

// 删除一条关系
func DeleteRelation(followerID int64, followedID int64) error {
	if followedID == followerID {
		return nil
	}
	return DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).Delete(&model.Relation{}).Error
}

// 查询一条关系是否存在
func IsFollow(followerID int64, followedID int64) (bool, error) {
	var relations []*model.Relation
	err := DB.Where("follower_id = ? AND followed_id =?", followerID, followedID).Find(&relations).Error
	if err != nil {
		return false, err
	}
	if len(relations) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// 获取关注
func GetFollowIDsByUserID(userID int64) ([]int64, error) {
	var relations []*model.Relation
	err := DB.Where("follower_id = ?", userID).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	var FollowedIDs []int64
	for _, fol := range relations {
		FollowedIDs = append(FollowedIDs, fol.FollowedID)
	}
	fmt.Printf("followed_id: %v\n", FollowedIDs)
	return FollowedIDs, nil
}

// 获取粉丝
func GetFanIDsByUserID(userID int64) ([]int64, error) {
	var relations []*model.Relation
	err := DB.Where("followed_id = ?", userID).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	var FollowerIDs []int64
	for _, fol := range relations {
		FollowerIDs = append(FollowerIDs, fol.FollowerID)
	}
	fmt.Printf("follower_id: %v\n", FollowerIDs)
	return FollowerIDs, nil
}

// 互关
func GetFriendIDsByUserID(userID int64) ([]int64, error) {
	followIDs, err := GetFollowIDsByUserID(userID)
	if err != nil {
		return nil, err
	}
	var friendIDs []int64
	for followID := range followIDs {
		var relations []*model.Relation
		err := DB.Where("follower_id = ? AND followed_id =?", followID, userID).Find(&relations).Error
		if err != nil {
			return nil, err
		}
		if len(relations) > 0 {
			friendIDs = append(friendIDs, int64(followID))
		}
	}
	return friendIDs, nil
}
