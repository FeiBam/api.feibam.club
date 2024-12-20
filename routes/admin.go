package routes

import (
	"api-feibam-club/middleware"
	"api-feibam-club/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(relativePath string, r *gin.Engine) {
	adminRoutes := r.Group(relativePath)

	adminRoutes.Use(middleware.IsLogin)

	adminRoutes.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, utils.JsonResponse("ok", 200, "You Are Success Login", "", gin.H{}))
	})
}
