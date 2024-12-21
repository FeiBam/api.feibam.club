package dao

import (
	"api-feibam-club/models"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

func CreateArticleModel(db *gorm.DB, articleData *models.ArticleFrontMatter) (*models.Article, error) {
	tags := []models.Tag{}
	for _, tagName := range articleData.Tags {
		var existingTag models.Tag
		if err := db.Where("name = ?", tagName).First(&existingTag).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to query tag: %w", err)
			}
			existingTag = models.Tag{Name: tagName}
			if err := db.Create(&existingTag).Error; err != nil {
				return nil, fmt.Errorf("failed to create tag: %w", err)
			}
		}
		tags = append(tags, existingTag)
	}

	links := []models.Link{}
	for _, linkURL := range articleData.Links {
		var existingLink models.Link
		if err := db.Where("url = ?", linkURL).First(&existingLink).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to query link: %w", err)
			}
			existingLink = models.Link{URL: linkURL}
			if err := db.Create(&existingLink).Error; err != nil {
				return nil, fmt.Errorf("failed to create link: %w", err)
			}
		}
		links = append(links, existingLink)
	}

	var langCode models.ArticleLangCode
	langCode, ok := langCode.FromString(articleData.Lang)
	if !ok {
		return nil, fmt.Errorf("unsupported language code")
	}

	createdAt, err := time.Parse(time.DateOnly, articleData.CreateAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CreateAt")
	}

	return &models.Article{
		ID:           uint(articleData.ID),
		Title:        articleData.Title,
		Introduction: articleData.Introduction,
		Tags:         tags,
		Subject:      articleData.Subject,
		CreatedAt:    createdAt,
		Lang:         langCode,
		Links:        links,
	}, nil
}

func InsertArticleRecord(db *gorm.DB, article *models.Article) error {
	tx := db.Begin()
	if err := tx.Create(article).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// updateArticleRecord 更新文章记录
func UpdateArticleRecord(db *gorm.DB, article *models.Article) error {
	var oldRecord models.Article

	// 预加载旧记录，包括关联数据
	err := db.Preload("Tags").Preload("Links").Where("id = ? AND lang = ?", article.ID, article.Lang).First(&oldRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No record found for article ID: %d", article.ID)
			return nil
		}
		return err
	}

	// 开始事务
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Transaction rolled back due to panic: %v", r)
		}
	}()

	// 更新主表记录
	if err := tx.Model(&oldRecord).Updates(article).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update article: %w", err)
	}

	// 确保关联记录存在并唯一
	for i, tag := range article.Tags {
		var existingTag models.Tag
		if err := tx.Where("name = ?", tag.Name).First(&existingTag).Error; err == nil {
			article.Tags[i] = existingTag // 使用已有记录
		}
	}
	for i, link := range article.Links {
		var existingLink models.Link
		if err := tx.Where("url = ?", link.URL).First(&existingLink).Error; err == nil {
			article.Links[i] = existingLink // 使用已有记录
		}
	}

	// 更新关联的 Tags 和 Links
	if err := tx.Model(&oldRecord).Association("Tags").Replace(article.Tags); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update tags: %w", err)
	}
	if err := tx.Model(&oldRecord).Association("Links").Replace(article.Links); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update links: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
