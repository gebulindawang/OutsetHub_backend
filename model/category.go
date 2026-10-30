package model

import "time"

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`       // 分类名,全局唯一
	Sort      int       `json:"sort"`       // 排序权重,越小越靠前
	CreatedAt time.Time `json:"created_at"`
}

func (c *Category) TableName() string {
	return "category"
}
