package model

import "time"

type Comment struct {
	ID            uint64    `json:"id"`
	ArticleID     uint64    `json:"article_id"`
	UserID        int       `json:"user_id"`
	ParentID      uint64    `json:"parent_id"`
	ReplyToUserID int       `json:"reply_to_user_id"`
	Content       string    `json:"content"`
	LikeCount     int       `json:"like_count"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

func (c *Comment) TableName() string {
	return "comment"
}
