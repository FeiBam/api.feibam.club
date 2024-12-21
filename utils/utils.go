package utils

import (
	"api-feibam-club/models"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
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

func ParseArticleData(data []byte, format string) (*models.ArticleFrontMatter, error) {
	var articleData models.ArticleFrontMatter

	switch format {
	case "yaml":
		if err := yaml.Unmarshal(data, &articleData); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	case "json":
		if err := json.Unmarshal(data, &articleData); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	default:
		return nil, errors.New("unsupported format: must be 'yaml' or 'json'")
	}

	return &articleData, nil
}

func ConvertCreateArticleJSONBindToFrontMatter(bind *models.CreateArticleJSONBind) *models.ArticleFrontMatter {
	return &models.ArticleFrontMatter{
		ID:           bind.ID,
		Title:        bind.Title,
		Introduction: bind.Introduction,
		Tags:         bind.Tags,
		CreateAt:     bind.CreateAt,
		Lang:         bind.Lang,
		Links:        bind.Links,
		Subject:      bind.Subject,
	}
}
func RespondWithError(ctx *gin.Context, statusCode int, message string, details interface{}) {
	ctx.JSON(statusCode, JsonResponse("err", statusCode, "", message, details))
}
