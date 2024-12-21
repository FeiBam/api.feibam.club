package routes

import (
	"api-feibam-club/controls"

	"github.com/gin-gonic/gin"
)

func ArticleRoutes(prefix string, group *gin.RouterGroup) {
	articleRoutes := group.Group(prefix)

	articleRoutes.GET("", controls.GetArticlesBy)

	articleRoutes.POST("", controls.CreateArticle)

	articleRoutes.GET("/:lang/:id", controls.GetArticleByLangWithId)

	articleRoutes.GET("/info", controls.InfoOfArticles)

}
