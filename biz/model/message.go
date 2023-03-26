package model

type Message struct {
	ID       int64  `json:"id"`
	Content  string `json:"content"`
	Type     int    `json:"type"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	SendTime int64  `json:"sendTime" bson:"sendTime"`
	Reached  bool   `json:"reached"`
}

type UserGroup struct {
	Username string `json:"username"`
	GroupID  int64  `json:"groupID" bson:"groupID"`
}

type GroupMessageReached struct {
	Receiver  string `json:"receiver"`
	GroupID   int64  `json:"groupID" bson:"groupID"`
	MessageID int64  `json:"messageID" bson:"messageID"`
}
