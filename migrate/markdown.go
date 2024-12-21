package migrate

import (
	"api-feibam-club/db"
	"api-feibam-club/db/dao"
	"api-feibam-club/models"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 错误定义
var (
	ErrReadMarkdownFile   = errors.New("failed to read Markdown file")
	ErrInvalidMarkdownFmt = errors.New("invalid Markdown format")
	ErrInvalidYAMLHeader  = errors.New("invalid YAML header")
)

// processMarkDownFile 解析 Markdown 文件
func processMarkDownFile(path string) (*models.ArticleFrontMatter, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: unable to read file %s", ErrReadMarkdownFile, path)
	}

	parts := strings.SplitN(string(fileContent), "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("%w: invalid file format in %s", ErrInvalidMarkdownFmt, path)
	}

	var articleData models.ArticleFrontMatter
	if err := yaml.Unmarshal([]byte(parts[1]), &articleData); err != nil {
		return nil, fmt.Errorf("%w: invalid YAML header in %s", ErrInvalidYAMLHeader, path)
	}
	articleData.Subject = parts[2]

	return &articleData, err
}

// MigrateMarkdowntoDatabases 启动迁移任务
func MigrateMarkdowntoDatabases() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		db, err := db.GetDB(logger.Silent)
		if err != nil {
			log.Fatalf("Failed to get database instance: %v", err)
		}

		dir, _ := cmd.Flags().GetString("path")
		force, _ := cmd.Flags().GetBool("force")
		update, _ := cmd.Flags().GetBool("update")
		if dir == "" {
			log.Println("You must specify the directory containing the Markdown files!")
			return
		}

		if update {
			log.Println("Starting multi-threaded update...")
			processUpdatesConcurrently(db, dir)
		} else {
			log.Println("Starting sequential creation of new articles...")
			processNewArticlesSequential(db, dir, force)
		}
	}
}

// 多线程处理更新
func processUpdatesConcurrently(db *gorm.DB, dir string) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // 控制并发数

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			articleData, err := processMarkDownFile(path)
			if err != nil {
				log.Printf("Failed to process file %s: %v", path, err)
				return
			}

			article, err := dao.CreateArticleModel(db, articleData)
			if err != nil {
				log.Printf("Failed to create article model for file %s: %v", path, err)
				return
			}

			if err := dao.UpdateArticleRecord(db, article); err != nil {
				log.Printf("Failed to update article for file %s: %v", path, err)
			} else {
				log.Printf("Updated article successfully. ID: %d, Path: %s", article.ID, path)
			}
		}()
		return nil
	})

	wg.Wait()
}

// 串行处理新建文章
func processNewArticlesSequential(db *gorm.DB, dir string, force bool) {
	var count int64
	if !force {
		if err := db.Model(&models.Article{}).Count(&count).Error; err != nil {
			log.Fatalf("Failed to count articles: %v", err)
		}
	}

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		log.Printf("Processing file: %s", path)
		articleData, err := processMarkDownFile(path)
		if err != nil {
			log.Printf("Failed to process file %s: %v", path, err)
			return nil
		}

		article, err := dao.CreateArticleModel(db, articleData)
		if err != nil {
			log.Printf("Failed to create article model for file %s: %v", path, err)
			return nil
		}

		if !force && article.ID != uint(count+1) {
			log.Printf("File ID mismatch for %s. Expected: %d, Got: %d", path, count+1, article.ID)
			return nil
		}

		if err := dao.InsertArticleRecord(db, article); err != nil {
			log.Printf("Failed to insert article for file %s: %v", path, err)
		} else {
			log.Printf("Inserted article successfully. ID: %d, Path: %s", article.ID, path)
			count++
		}
		return nil
	})
}

// insertArticleRecord 插入文章记录
