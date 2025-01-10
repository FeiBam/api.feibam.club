package routes

import (
	"api-feibam-club/controls"
	"api-feibam-club/middleware"
	"api-feibam-club/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(prefix string, group *gin.RouterGroup) {
	adminRoutes := group.Group(prefix)

	adminRoutes.POST("/login", controls.Login)

	adminRoutes.Use(middleware.IsLogin)

	adminRoutes.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, utils.JsonResponse("ok", 200, "You Are./ Success Login", "", gin.H{}))
	})
}
