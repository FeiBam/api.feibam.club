package routes

import (
	"api-feibam-club/controls"
	"api-feibam-club/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ArticleRoutes(prefix string, group *gin.RouterGroup) {
	articleRoutes := group.Group(prefix)

	articleRoutes.GET("", controls.GetArticlesBy)

	articleRoutes.GET("/:lang/:id", controls.GetArticleByLangWithId)

	articleRoutes.GET("/info", controls.InfoOfArticles)

	utils.RegisterRoutes("/teapot", articleRoutes, i_am_teapot)
}

func i_am_teapot(prefix string, group *gin.RouterGroup) {
	i_am_teapot_routes := group.Group(prefix)

	i_am_teapot_routes.Any("", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "https://www.google.com/teapot")

	})
}
