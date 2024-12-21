package controls

import (
	"api-feibam-club/db/dao"
	"api-feibam-club/models"
	"api-feibam-club/utils"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InfoOfArticles(ctx *gin.Context) {
	var db *gorm.DB
	var err error
	var totalArticleCount int64
	langCounts := gin.H{}
	tagCounts := gin.H{}
	var tags []models.Tag
	var langCode models.ArticleLangCode

	var queryParament models.GetArticleInfoByFormBind
	if err := ctx.ShouldBindQuery(&queryParament); err != nil {
		ctx.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	// 获取数据库实例
	db, err = utils.GetDBFromContext(ctx)
	if err != nil {
		ctx.JSON(500, utils.JsonResponse("err", 500, "", err.Error(), gin.H{}))
		return
	}

	// 获取总文章数量
	if err := db.Model(&models.Article{}).Count(&totalArticleCount).Error; err != nil {
		ctx.JSON(500, utils.JsonResponse("err", 500, "", err.Error(), gin.H{}))
		return
	}

	// 获取每种语言的文章数量
	iter := langCode.GetStringIterators()
	for {
		lang, ok := iter()
		if !ok {
			break
		}
		codeOfLang, ok := langCode.FromString(lang)
		if !ok {
			continue
		}
		var count int64
		if err := db.Model(&models.Article{}).Where("lang = ?", codeOfLang).Count(&count).Error; err != nil {
			fmt.Printf("Failed to count articles for lang %v: %v\n", lang, err)
			langCounts[lang] = 0 // 出错时设为 0
			continue
		}
		langCounts[lang] = count
	}

	// 获取所有标签
	if err := db.Select("id, name").Find(&tags).Error; err != nil {
		ctx.JSON(500, utils.JsonResponse("err", 500, "", err.Error(), gin.H{}))
		return
	}

	langCode, ok := langCode.FromString(queryParament.Lang)

	if !ok {
		ctx.JSON(400, utils.JsonResponse("ok", 200, "", "Unsupported language codes.", gin.H{}))
		return
	}
	// 统计每个标签关联的文章数量
	for _, tag := range tags {
		var count int64
		if err := db.Table("articles").
			Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ? AND articles.lang = ?", tag.ID, langCode).
			Count(&count).Error; err != nil {
			fmt.Printf("Failed to count articles for tag %v: %v\n", tag.Name, err)
			tagCounts[tag.Name] = 0 // 出错时设为 0
			continue
		}
		tagCounts[tag.Name] = count
	}

	// 返回结果
	ctx.JSON(200, utils.JsonResponse("ok", 200, "", "", gin.H{
		"articleCount":             totalArticleCount,
		"articleCountOfLang":       langCounts,
		"articleCountOfLangAndTag": tagCounts,
	}))
}

func GetArticlesByLang(ctx *gin.Context, size int, page int, lang_code models.ArticleLangCode) {
	var db *gorm.DB
	var articles []models.Article
	var err error

	db, err = utils.GetDBFromContext(ctx)
	if err != nil {
		ctx.JSON(500, utils.JsonResponse("err", 500, "", err.Error(), gin.H{}))
		return
	}
	result := db.Limit(size).Preload("Tags").Preload("Links").Offset((page-1)*size).Where("lang = ?", lang_code).Find(&articles)
	if result.Error != nil {
		ctx.JSON(500, utils.JsonResponse("err", 500, "", result.Error.Error(), gin.H{}))
		return
	}

	var articleDTOs []models.ArticleDTO
	for _, article := range articles {
		articleDTOs = append(articleDTOs, utils.ToArticleDTO(article))
	}

	ctx.JSON(200, utils.JsonResponse("ok", 200, "", "", gin.H{
		"articles": articleDTOs,
	}))
}

func GetArticleByTagWithLang(ctx *gin.Context, size int, page int, lang_code models.ArticleLangCode, tag string) {
	var db *gorm.DB
	var articles []models.Article
	var err error
	db, err = utils.GetDBFromContext(ctx)
	if err != nil {
		ctx.JSON(500, utils.JsonResponse("err", 500, "", err.Error(), gin.H{}))
		return
	}
	result := db.Limit(size).
		Preload("Tags").
		Preload("Links").
		Offset((page-1)*size).
		Joins("JOIN article_tags ON article_tags.article_id = articles.id").
		Joins("JOIN tags ON article_tags.tag_id = tags.id").
		Where("tags.name = ? AND articles.lang = ?", tag, lang_code).
		Find(&articles)
	if result.Error != nil {
		ctx.JSON(500, utils.JsonResponse("err", 500, "", result.Error.Error(), gin.H{}))
		return
	}

	var articleDTOs []models.ArticleDTO
	for _, article := range articles {
		articleDTOs = append(articleDTOs, utils.ToArticleDTO(article))
	}
	ctx.JSON(200, utils.JsonResponse("ok", 200, "", "", gin.H{
		"articles": articleDTOs,
	}))
}

func GetArticlesBy(ctx *gin.Context) {
	var ok bool
	var queryParament models.GetArticlesByFormBind
	if err := ctx.BindQuery(&queryParament); err != nil {
		ctx.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	var lang_code models.ArticleLangCode
	lang_code, ok = lang_code.FromString(queryParament.Lang)
	if !ok {
		ctx.JSON(400, utils.JsonResponse("ok", 200, "", "Unsupported language codes.", gin.H{}))
		return
	}
	fmt.Print(queryParament.Tag)
	if queryParament.Tag != "" {
		fmt.Print("use with tag")
		GetArticleByTagWithLang(ctx, queryParament.Size, queryParament.Page, lang_code, queryParament.Tag)
		return
	}
	GetArticlesByLang(ctx, queryParament.Size, queryParament.Page, lang_code)
}

func GetArticleByLangWithId(ctx *gin.Context) {
	var db *gorm.DB
	var ok bool
	var err error
	var lang_code models.ArticleLangCode
	var article models.Article
	var queryParament models.GetArticleByLangWithIdUrlBind

	if err := ctx.ShouldBindUri(&queryParament); err != nil {
		ctx.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	lang_code, ok = lang_code.FromString(queryParament.Lang)

	if !ok {
		ctx.JSON(400, utils.JsonResponse("ok", 200, "", "Unsupported language codes.", gin.H{}))
		return
	}
	db, err = utils.GetDBFromContext(ctx)
	if err != nil {
		ctx.JSON(500, utils.JsonResponse("err", 500, "", err.Error(), gin.H{}))
		return
	}

	err = db.Preload("Tags").Preload("Links").Where("id = ? AND lang = ?", queryParament.Id, lang_code).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(404, utils.JsonResponse("ok", 200, "", "No article records matching this language and this id were found.", gin.H{}))
			return
		}
		ctx.JSON(500, utils.JsonResponse("err", 500, "", fmt.Sprintf("An unexpected error occurred! err:%v", err), gin.H{}))
		return
	}
	articleDTO := utils.ToArticleDTO(article)
	ctx.JSON(200, utils.JsonResponse("ok", 200, "", "", articleDTO))
}

func CreateArticle(ctx *gin.Context) {
	var count int64
	var db *gorm.DB
	var err error
	var queryParament models.CreateArticleJSONBind

	// 解析请求体
	if err := ctx.ShouldBindJSON(&queryParament); err != nil {
		if err.Error() == "EOF" {
			utils.RespondWithError(ctx, 400, "The articleData field is required!", nil)
			return
		}
		utils.RespondWithError(ctx, 400, err.Error(), nil)
		return
	}

	// 获取数据库连接
	db, err = utils.GetDBFromContext(ctx)
	if err != nil {
		utils.RespondWithError(ctx, 500, "Failed to get database connection", nil)
		return
	}

	// 获取当前文章计数
	if err := db.Model(&models.Article{}).Count(&count).Error; err != nil {
		utils.RespondWithError(ctx, 500, "Failed to count articles", nil)
		return
	}

	// 转换为 FrontMatter
	articleData := utils.ConvertCreateArticleJSONBindToFrontMatter(&queryParament)

	// 校验 Article ID
	expectedID := uint(count + 1)
	if !validateArticleID(ctx, articleData.ID, expectedID) {
		return
	}
	articleData.ID = int(expectedID)

	// 创建文章
	article, err := dao.CreateArticleModel(db, articleData)
	if err != nil {
		utils.RespondWithError(ctx, 400, "Failed to create article", err.Error())
		return
	}

	// 插入文章记录
	if err := dao.InsertArticleRecord(db, article); err != nil {
		utils.RespondWithError(ctx, 500, "Failed to insert article record", err.Error())
		return
	}

	// 转换为 DTO 并返回成功响应
	articleDTO := utils.ToArticleDTO(*article)
	ctx.JSON(200, utils.JsonResponse("ok", 200, "Success Create Article!", "", articleDTO))
}

func validateArticleID(ctx *gin.Context, articleID int, expectedID uint) bool {
	if articleID != -1 && uint(articleID) != expectedID {
		errMsg := fmt.Sprintf("Invalid Article ID! Expected %d but got %d", expectedID, articleID)
		utils.RespondWithError(ctx, 400, errMsg, nil)
		return false
	}
	return true
}
