package dao

import "backend/model"

func GetArticles() ([]model.Article, error) {
	var articles []model.Article
	result := DB.Where("status = ?", 1).Order("created_at DESC").Limit(10).Find(&articles)
	return articles, result.Error
}

func GetArticleByID(id uint64) (model.Article, error) {
	var article model.Article
	err := DB.First(&article, id).Error
	return article, err
}

func CreateArticle(article model.Article) error {
	return DB.Create(&article).Error
}

func UpdateArticle(id uint64, updates map[string]any) error {
	return DB.Model(&model.Article{}).Where("id = ?", id).Updates(updates).Error
}

func DeleteArticle(id uint64) error {
	return DB.Model(&model.Article{}).Where("id = ?", id).Update("status", 2).Error
}
