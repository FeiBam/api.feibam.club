package main

import (
	"api-feibam-club/db"
	"api-feibam-club/middleware"
	"api-feibam-club/migrate"
	"api-feibam-club/models"
	"api-feibam-club/routes"
	"api-feibam-club/utils"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm/logger"
)

func addDB(ctx *gin.Context) {
	db, err := db.GetDB(logger.Silent)
	if err != nil {
		panic(fmt.Sprintf("fail get databases.... error : %v", err))
	}
	ctx.Set("db", db)
	ctx.Next()
}

func runServer() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		if port == "" {
			port = ":8080"
		}

		fmt.Printf("Running server on port %s\n", port)

		r := gin.Default()

		r.TrustedPlatform = gin.PlatformCloudflare

		r.Use(middleware.XResponseTime)

		r.Use(middleware.SecurityHeaders)

		r.Use(addDB)

		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"}, // 必须明确指定来源
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: false, // 如果请求携带凭证（如 Cookies）
			MaxAge:           0,     // 缓存时间
		}))

		defalut_route := r.Group("")
		defalut_route.Any("/teapot", func(ctx *gin.Context) {
			ctx.Redirect(http.StatusTemporaryRedirect, "https://www.google.com/teapot")
		})

		utils.RegisterRoutes("/article", defalut_route, routes.ArticleRoutes)

		if err := r.Run(port); err != nil {
			panic(fmt.Sprintf("failed to start server: %v", err))
		}
	}
}

func migrateModeltoDatabases() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		db, err := db.GetDB(logger.Silent)
		if err != nil {
			panic(fmt.Sprintf("failed to get database instance... error: %v", err))
		}
		fmt.Println("Migrating models to the database...")
		db.AutoMigrate(&models.Article{})
		db.AutoMigrate(&models.Link{})
		db.AutoMigrate(&models.Tag{})
		fmt.Println("Migration completed.")
	}
}

func main() {
	var rootCmd = &cobra.Command{}

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a command",
	}

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		Run:   runServer(),
	}
	serverCmd.Flags().StringP("port", "p", ":8080", "Port to run the server on")

	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate model or markdown to databases",
	}

	var migrateModelCmd = &cobra.Command{
		Use:   "db",
		Short: "Migrate Model to databases",
		Run:   migrateModeltoDatabases(),
	}

	var migrateMarkdownCmd = &cobra.Command{
		Use:   "md",
		Short: "Migrate markdown file to databases",
		Run:   migrate.MigrateMarkdowntoDatabases(),
	}
	migrateMarkdownCmd.Flags().StringP("path", "p", "", "Markdown file directory location")
	migrateMarkdownCmd.Flags().Bool("force", false, "Force overwrite database record")
	migrateMarkdownCmd.Flags().Bool("update", false, "Update database record")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(migrateCmd)

	runCmd.AddCommand(serverCmd)

	migrateCmd.AddCommand(migrateMarkdownCmd)
	migrateCmd.AddCommand(migrateModelCmd)

	rootCmd.Execute()
}
