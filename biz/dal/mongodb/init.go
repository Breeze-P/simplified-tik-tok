package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client          *mongo.Client
	err             error
	db              *mongo.Database
	singleMsg       *mongo.Collection
	groupMsg        *mongo.Collection
	groupMsgReached *mongo.Collection
	userGroup       *mongo.Collection
)

// Init 连接远程mongodb
func Init() {
	ClientOpts := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s",
			os.Getenv("MONGO_HOST"),
			os.Getenv("MONGO_PORT")),
		).
		SetAuth(options.Credential{
			Username: os.Getenv("MONGO_USER"),
			Password: os.Getenv("MONGO_PASSWORD"),
		}).
		SetConnectTimeout(10 * time.Second)

	// 1.建立连接
	if client, err = mongo.Connect(context.TODO(), ClientOpts); err != nil {
		fmt.Print(err)
		return
	}
	// 2.选择数据库
	db = client.Database("msg")
	collection, err := db.ListCollectionNames(context.TODO(), bson.M{})
	log.Println(collection, err)
	// 3.选择表
	singleMsg = db.Collection("single_msg")
	groupMsg = db.Collection("group_msg")
	groupMsgReached = db.Collection("group_msg_reached")
	userGroup = db.Collection("user_group")
}
