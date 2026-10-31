package model

import "time"

type Article struct {
	ID           uint64    `json:"id"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Content      string    `json:"content"`
	UserID       int       `json:"user_id"`
	CategoryID   int       `json:"category_id"`
	ViewCount    int       `json:"view_count"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (a *Article) TableName() string {
	return "article"
}
