package dao

import "backend/model"

func GetComments(articleID uint64) ([]model.Comment, error) {
	var comments []model.Comment
	result := DB.Where("article_id = ? AND status = ?", articleID, 1).Order("created_at ASC").Find(&comments)
	return comments, result.Error
}

func GetCommentByID(id uint64) (model.Comment, error) {
	var comment model.Comment
	err := DB.First(&comment, id).Error
	return comment, err
}

func CreateComment(comment model.Comment) error {
	return DB.Create(&comment).Error
}

func UpdateComment(id uint64, updates map[string]any) error {
	return DB.Model(&model.Comment{}).Where("id = ?", id).Updates(updates).Error
}

func DeleteComment(id uint64) error {
	return DB.Model(&model.Comment{}).Where("id = ?", id).Update("status", 2).Error
}
