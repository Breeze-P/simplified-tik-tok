package resbody

type Comment struct {
	ID         int64  `json:"id"`          // 视频评论id
	User       User   `json:"user"`        // 评论用户信息
	Content    string `json:"content"`     // 评论内容
	CreateDate string `json:"create_data"` // 评论发布日期，格式 mm-dd
}
