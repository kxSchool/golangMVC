package ado

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"golang.com/models"
	"golang.com/utils"
)

// NewsIndex is yes.
func NewsIndex(db *sql.DB) (*models.Pictures, error) {
	list, err := db.Query("select ID,Subject,Picture from Article where PostType=? order by AddDate desc limit ?", "新闻", 30)
	utils.CheckErr(err)

	pictures := models.Pictures{}
	pictures.Items = []*models.Picture{}
	for list.Next() {
		picture := models.Picture{}
		err = list.Scan(&picture.ID, &picture.Subject, &picture.URL)
		utils.CheckErr(err)

		pictures.Items = append(pictures.Items, &picture)
	}

	list.Close()
	return &pictures, nil
}

// TechIndex is yes.
func TechIndex(db *sql.DB) (*models.Course, error) {
	list, err := db.Query("select ID,Subject,Picture,Description from Article order by AddDate desc limit ?", 30)
	utils.CheckErr(err)
	course := models.Course{}
	course.ArticleItems = []*models.Article{}
	for list.Next() {
		var article models.Article
		err = list.Scan(&article.ID, &article.Subject, &article.Picture, &article.Description)
		utils.CheckErr(err)

		course.ArticleItems = append(course.ArticleItems, &article)
	}

	list.Close()
	return &course, nil
}
