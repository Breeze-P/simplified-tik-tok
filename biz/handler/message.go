package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	hzutils "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/websocket"

	"simplified-tik-tok/biz/dal/mongodb"
	"simplified-tik-tok/biz/model"
	"simplified-tik-tok/biz/utils"
)

const (
	// 私聊消息
	privateMsg = 1
	// 群聊消息
	groupMsg = 2
)

const (
	writeChannelBuffer = 3
)

type connectStruct struct {
	UserID string `json:"userId"`
}

type responseStruct struct {
	Status bool `json:"status"`
}

// websocket upgrader
var upgrader = websocket.HertzUpgrader{} // use default options

var clientSenders = make(map[string]chan model.Message)

// Sender2Client 负责从待发送消息队列中获取消息，发送给客户端
func Sender2Client(username string, conn *websocket.Conn) {
	for {
		select {
		case msg := <-clientSenders[username]:
			if err := conn.WriteJSON(&msg); err != nil {
				log.Println("Sender2Client write:", err)
				// 私聊不在线存入数据库，群聊不用管
				if msg.Type == privateMsg {
					mongodb.InsertSingleMsg(msg)
				}
				return
			}
			// 私聊在线存入数据库
			if msg.Type == privateMsg {
				msg.Reached = true
				mongodb.InsertSingleMsg(msg)
			} else if msg.Type == groupMsg {
				// 接收到再单独插入GroupMessageReachedStruct表
				groupID, _ := strconv.Atoi(msg.Receiver)
				mongodb.InsertGroupMessageReached(model.GroupMessageReached{
					Receiver:  username,
					GroupID:   int64(groupID),
					MessageID: msg.ID,
				})
			}
		}
	}
}

// Message
func Message(_ context.Context, c *app.RequestContext) {
	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		// 获取用户ID
		var userInfo connectStruct
		err := conn.ReadJSON(&userInfo)
		userID := userInfo.UserID
		if err != nil {
			log.Println("read userid:", err)
			return
		} else if userID == "" {
			log.Println("userInfo format error")
			return
		}
		log.Printf("%s online.", userID)
		// 返回未到达的消息记录
		// 查询未到达的私聊消息
		singleUnreachedMsgs := mongodb.GetSingleMsg(userID, false)
		// 查询未到达的群聊消息
		groupUnreachedMsgs := mongodb.GetGroupMsg(userID, false)
		// err = conn.WriteJSON(&responseStruct{Status: true})
		log.Println("singleUnreachedMsgs:", singleUnreachedMsgs)
		log.Println("groupUnreachedMsgs:", groupUnreachedMsgs)
		if err = conn.WriteJSON(&singleUnreachedMsgs); err != nil {
			log.Println("write single:", err)
			return
		}
		if err = conn.WriteJSON(&groupUnreachedMsgs); err != nil {
			log.Println("write group:", err)
			return
		}
		// 消息增量抵达成功
		mongodb.UpdateSingleMsgReached(singleUnreachedMsgs)
		mongodb.UpdateGroupMsgReached(groupUnreachedMsgs, userID)
		// 类似于一个消息队列，用于存放待发送的消息
		writeChannel := make(chan model.Message, writeChannelBuffer)
		clientSenders[userID] = writeChannel
		// 启动协程：从消息队列中获取消息，发送给客户端
		go Sender2Client(userID, conn)

		// 后续负责读取客户端发送的消息
		for {
			var msg model.Message
			if err = conn.ReadJSON(&msg); err != nil {
				log.Println("for read:", err)
				delete(clientSenders, userID)
				break
			}
			// 如果消息不带ID，需要生成ID
			if msg.ID == 0 {
				msg.ID = utils.GenerateSnowflake()
			}
			if msg.Type == privateMsg {
				receiver := msg.Receiver
				// 查看接收者是否在线，在线则交给接收者的消息队列，否则存入数据库
				if _, ok := clientSenders[receiver]; ok {
					clientSenders[receiver] <- msg
				} else {
					// 不在线存入数据库
					mongodb.InsertSingleMsg(msg)
				}
			} else if msg.Type == groupMsg {
				// 群聊设计与私聊不同，不考虑Reached字段，是否接收到由GroupMessageReachedStruct表记录
				groupIDStr := msg.Receiver
				groupID, err := strconv.Atoi(groupIDStr)
				if err != nil {
					log.Println("group id is not int")
					continue
				}
				// 因此统一存入数据库，接收到再单独插入GroupMessageReachedStruct表
				mongodb.InsertGroupMsg(msg)
				log.Println("group id:", groupID)
				receivers := mongodb.GetGroupUsers(int64(groupID))
				for _, receiver := range receivers {
					// log.Println(receiver, "in turn")
					if _, ok := clientSenders[receiver]; ok {
						// log.Println(receiver, "in sender")
						clientSenders[receiver] <- msg
					}
				}
			}

		}
	})
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
}

func MessageHistory(_ context.Context, c *app.RequestContext) {
	var userInfo connectStruct
	if err := c.BindAndValidate(&userInfo); err != nil {
		c.JSON(http.StatusOK, hzutils.H{
			"status_code": http.StatusBadRequest,
			"status_msg":  err.Error(),
			"single_msg":  nil,
			"group_msg":   nil,
		})
		return
	}
	userID := userInfo.UserID
	c.JSON(http.StatusOK, hzutils.H{
		"status_code": http.StatusOK,
		"status_msg":  "success",
		"single_msg":  mongodb.GetSingleMsg(userID, true),
		"group_msg":   mongodb.GetGroupMsg(userID, true),
	})
}
