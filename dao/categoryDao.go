package dao

import "backend/model"

func GetCategorys() ([]model.Category, error) {
	var categorys []model.Category
	result := DB.Order("sort ASC, id ASC").Find(&categorys)
	return categorys, result.Error
}

func CreateCategory(category model.Category) error {
	return DB.Create(&category).Error
}

func UpdateCategory(id int, updates map[string]any) error {
	return DB.Model(&model.Category{}).Where("id = ?", id).Updates(updates).Error
}

func DeleteCategory(id int) error {
	return DB.Delete(&model.Category{}, id).Error
}
