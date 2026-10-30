package dao

import "backend/model"

func GetCategorys()([]model.Category,error) {
	var categorys []model.Category
	err := DB.Find(&categorys)
	return categorys,err.Error
}