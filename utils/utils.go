package utils

import (
	"api-feibam-club/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JsonResponse(status string, code int, message string, err string, body any) gin.H {
	return gin.H{
		"status":  status,
		"code":    code,
		"error":   err,
		"message": message,
		"body":    body,
	}
}

func RegisterRoutes(relativePath string, target *gin.RouterGroup, source func(prefix string, target *gin.RouterGroup)) {
	source(relativePath, target)
}

func GetDBFromContext(ctx *gin.Context) (*gorm.DB, error) {
	value, exists := ctx.Get("db")
	if !exists {
		return nil, fmt.Errorf("database instance not found in context")
	}
	db, ok := value.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("failed to assert database instance")
	}
	return db, nil
}

func ToArticleDTO(article models.Article) models.ArticleDTO {
	tags := make([]models.TagDTO, len(article.Tags))
	for i, tag := range article.Tags {
		tags[i] = models.TagDTO{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	links := make([]models.LinkDTO, len(article.Links))
	for i, link := range article.Links {
		links[i] = models.LinkDTO{
			URL: link.URL,
		}
	}

	return models.ArticleDTO{
		ID:           article.ID,
		Title:        article.Title,
		Introduction: article.Introduction,
		CreateAt:     article.CreatedAt.Format(time.DateOnly),
		Subject:      article.Subject,
		Lang:         int(article.Lang),
		Tags:         tags,
		Links:        links,
	}
}
