package model

import "time"

type Article struct {
	ID           uint64
	Title        string
	Summary      string
	Content      string
	UserID       int
	CategoryID   int
	ViewCount    int
	LikeCount    int
	CommentCount int
	Status       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (a *Article) TableName() string {
	return "article"
}
