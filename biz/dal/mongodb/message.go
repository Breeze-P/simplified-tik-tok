package mongodb

import (
	"context"
	"strconv"

	"simplified-tik-tok/biz/model"

	"go.mongodb.org/mongo-driver/bson"
)

// InsertSingleMsg 插入一条私聊数据
func InsertSingleMsg(msg model.Message) {
	_, err := singleMsg.InsertOne(context.TODO(), msg)
	if err != nil {
		return
	}
}

// InsertGroupMsg 插入一条群聊数据
func InsertGroupMsg(msg model.Message) {
	_, err := groupMsg.InsertOne(context.TODO(), msg)
	if err != nil {
		return
	}
}

// GetSingleMsg 获取用户的私聊消息
// @param username 用户名
// @param history 是否查询历史数据, true: 历史全部发送数据以及到达自己的数据，false: 未到达自己的数据
func GetSingleMsg(username string, history bool) []model.Message {
	var msg []model.Message
	filter := bson.M{"receiver": username, "reached": history}
	cur, err := singleMsg.Find(context.TODO(), filter)
	if err != nil {
		return msg
	}
	for cur.Next(context.TODO()) {
		var tmp model.Message
		err := cur.Decode(&tmp)
		if err != nil {
			return msg
		}
		msg = append(msg, tmp)
	}
	// 如果查询历史数据，还要包含发送的全部数据
	if history {
		filter = bson.M{"sender": username}
		cur, err = singleMsg.Find(context.TODO(), filter)
		if err != nil {
			return msg
		}
		for cur.Next(context.TODO()) {
			var tmp model.Message
			err := cur.Decode(&tmp)
			if err != nil {
				return msg
			}
			msg = append(msg, tmp)
		}
	}
	return msg
}

// GetGroupUsers 获取某个群聊的所有用户
func GetGroupUsers(groupID int64) []string {
	var users []string
	filter := bson.M{"groupID": groupID}
	cur, err := userGroup.Find(context.TODO(), filter)
	if err != nil {
		return users
	}
	for cur.Next(context.TODO()) {
		var tmp model.UserGroup
		err := cur.Decode(&tmp)
		if err != nil {
			return users
		}
		users = append(users, tmp.Username)
	}
	return users
}

// GetGroupMsg 获取某个人的群聊消息
// @param reached: true表示获取已经到达username的消息（群组特性包含自己发送的数据），false表示获取未到达username的消息
// @return: []model.MessageStruct
func GetGroupMsg(username string, reached bool) []model.Message {
	// 获取用户已经在GroupMessageReachedStruct中的记录
	var gmr []int64
	filter := bson.M{"receiver": username}
	cur, err := groupMsgReached.Find(context.TODO(), filter)
	if err != nil {
		return nil
	}
	for cur.Next(context.TODO()) {
		var tmp model.GroupMessageReached
		err := cur.Decode(&tmp)
		if err != nil {
			return nil
		}
		gmr = append(gmr, tmp.MessageID)
	}
	// 获取这个用户所有群聊
	var groups []string
	filter = bson.M{"username": username}
	cur, err = userGroup.Find(context.TODO(), filter)
	if err != nil {
		return nil
	}
	for cur.Next(context.TODO()) {
		var tmp model.UserGroup
		err := cur.Decode(&tmp)
		if err != nil {
			return nil
		}
		groups = append(groups, strconv.Itoa(int(tmp.GroupID)))
	}
	var msg []model.Message
	if reached {
		filter = bson.M{"receiver": bson.M{"$in": groups}, "id": bson.M{"$in": gmr}}
	} else {
		//{"receiver":{"$in":["1"]},"id":{"$nin":[1618519172319285248,1618520157112504320]}}
		filter = bson.M{"receiver": bson.M{"$in": groups}, "id": bson.M{"$nin": gmr}}
	}
	cur, err = groupMsg.Find(context.TODO(), filter)
	if err != nil {
		return msg
	}
	for cur.Next(context.TODO()) {
		var tmp model.Message
		err := cur.Decode(&tmp)
		if err != nil {
			return msg
		}
		msg = append(msg, tmp)
	}
	return msg
}

// UpdateSingleMsgReached 将指定的私聊消息标记为已到达
func UpdateSingleMsgReached(msgs []model.Message) {
	for _, msg := range msgs {
		filter := bson.M{"id": msg.ID}
		update := bson.M{"$set": bson.M{"reached": true}}
		_, err2 := singleMsg.UpdateOne(context.TODO(), filter, update)
		if err2 != nil {
			return
		}
	}
}

// UpdateGroupMsgReached 将指定的群聊消息标记为已到达
func UpdateGroupMsgReached(msgs []model.Message, username string) {
	for _, msg := range msgs {
		groupid, _ := strconv.Atoi(msg.Receiver)
		InsertGroupMessageReached(model.GroupMessageReached{
			Receiver:  username,
			MessageID: msg.ID,
			GroupID:   int64(groupid),
		})
	}
}

// InsertGroupMessageReached 插入一条群聊消息已到达记录
func InsertGroupMessageReached(gmr model.GroupMessageReached) {
	_, err := groupMsgReached.InsertOne(context.TODO(), gmr)
	if err != nil {
		return
	}
}
