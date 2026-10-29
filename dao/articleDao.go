package dao

import (
	"backend/model"
)

func GetArticles()([]model.Article,error){
	var articles []model.Article
	res := DB.Limit(10).Find(&articles)
	return articles,res.Error
}

func CreateArticle(article model.Article) error{
	res :=DB.Create(&article)
	return res.Error
}

func DeleteArticle(article model.Article) error{
	res := DB.Where("id = ?" ,article.ID).Delete(&article)
	return res.Error
}

func UpdateArticle(article model.Article)error{
	res := DB.Model(&article).Where("id = ?",article.ID).Updates(&article)
	return  res.Error
}