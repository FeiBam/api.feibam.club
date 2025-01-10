package controls

import (
	"api-feibam-club/models"
	"api-feibam-club/utils"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Login(ctx *gin.Context) {
	var json_parament models.AdminLoginJsonBind
	if err := ctx.BindJSON(&json_parament); err != nil {
		ctx.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	username := os.Getenv("USER_NAME")
	password := os.Getenv("PASS_WORD")
	fmt.Println(username, "|", json_parament.UserName)
	fmt.Println(password, "|", json_parament.PassWord)

	if json_parament.UserName != username || json_parament.PassWord != password {
		ctx.JSON(401, utils.JsonResponse("error", 401, "Incorrect username or password", "", nil))
		return
	}

	exp := time.Now().Add(time.Hour).Unix()

	jwtClaims := &jwt.MapClaims{
		"user_name": username,
		"exp":       exp,
	}

	token, err := utils.GenerateJWT(jwtClaims)
	if err != nil {
		ctx.JSON(500, utils.JsonResponse("error", 500, "Failed to generate token", "", nil))
		return
	}
	tokenInfo := &models.TokenInfo{
		Token:     token,
		ExpiresAt: exp,
	}

	// 获取 tokenStore
	tokenStore := utils.GetTokenStoreFromContext(ctx)

	// 将 token 存入 tokenStore
	tokenStore.Add(json_parament.UserName, *tokenInfo)

	// 返回成功响应
	ctx.JSON(200, utils.JsonResponse("ok", 200, "Login successful", "", tokenInfo))
}
