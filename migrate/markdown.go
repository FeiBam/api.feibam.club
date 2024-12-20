package migrate

import (
	"api-feibam-club/db"
	"api-feibam-club/models"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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
func processMarkDownFile(path string) (*models.ArticleMarkdownData, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: unable to read file %s", ErrReadMarkdownFile, path)
	}

	parts := strings.SplitN(string(fileContent), "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("%w: invalid file format in %s", ErrInvalidMarkdownFmt, path)
	}

	var metadata models.FrontMatter
	if err := yaml.Unmarshal([]byte(parts[1]), &metadata); err != nil {
		return nil, fmt.Errorf("%w: invalid YAML header in %s", ErrInvalidYAMLHeader, path)
	}

	return &models.ArticleMarkdownData{
		MetaData: metadata,
		Subject:  parts[2],
	}, nil
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

			markdownData, err := processMarkDownFile(path)
			if err != nil {
				log.Printf("Failed to process file %s: %v", path, err)
				return
			}

			article, err := createArticleModel(db, path, markdownData)
			if err != nil {
				log.Printf("Failed to create article model for file %s: %v", path, err)
				return
			}

			if err := updateArticleRecord(db, article); err != nil {
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
		markdownData, err := processMarkDownFile(path)
		if err != nil {
			log.Printf("Failed to process file %s: %v", path, err)
			return nil
		}

		article, err := createArticleModel(db, path, markdownData)
		if err != nil {
			log.Printf("Failed to create article model for file %s: %v", path, err)
			return nil
		}

		if !force && article.ID != uint(count+1) {
			log.Printf("File ID mismatch for %s. Expected: %d, Got: %d", path, count+1, article.ID)
			return nil
		}

		if err := insertArticleRecord(db, article); err != nil {
			log.Printf("Failed to insert article for file %s: %v", path, err)
		} else {
			log.Printf("Inserted article successfully. ID: %d, Path: %s", article.ID, path)
			count++
		}
		return nil
	})
}

// insertArticleRecord 插入文章记录
func insertArticleRecord(db *gorm.DB, article *models.Article) error {
	tx := db.Begin()
	if err := tx.Create(article).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// updateArticleRecord 更新文章记录
func updateArticleRecord(db *gorm.DB, article *models.Article) error {
	var oldRecord models.Article

	// 预加载旧记录，包括关联数据
	err := db.Preload("Tags").Preload("Links").Where("id = ? AND lang = ?", article.ID, article.Lang).First(&oldRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No record found for article ID: %d", article.ID)
			return nil
		}
		return err
	}

	// 开始事务
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Transaction rolled back due to panic: %v", r)
		}
	}()

	// 更新主表记录
	if err := tx.Model(&oldRecord).Updates(article).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update article: %w", err)
	}

	// 确保关联记录存在并唯一
	for i, tag := range article.Tags {
		var existingTag models.Tag
		if err := tx.Where("name = ?", tag.Name).First(&existingTag).Error; err == nil {
			article.Tags[i] = existingTag // 使用已有记录
		}
	}
	for i, link := range article.Links {
		var existingLink models.Link
		if err := tx.Where("url = ?", link.URL).First(&existingLink).Error; err == nil {
			article.Links[i] = existingLink // 使用已有记录
		}
	}

	// 更新关联的 Tags 和 Links
	if err := tx.Model(&oldRecord).Association("Tags").Replace(article.Tags); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update tags: %w", err)
	}
	if err := tx.Model(&oldRecord).Association("Links").Replace(article.Links); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update links: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// createArticleModel 创建文章模型
func createArticleModel(db *gorm.DB, path string, markdownData *models.ArticleMarkdownData) (*models.Article, error) {
	tags := []models.Tag{}
	for _, tagName := range markdownData.MetaData.Tags {
		var existingTag models.Tag
		if err := db.Where("name = ?", tagName).First(&existingTag).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to query tag: %w", err)
			}
			existingTag = models.Tag{Name: tagName}
			if err := db.Create(&existingTag).Error; err != nil {
				return nil, fmt.Errorf("failed to create tag: %w", err)
			}
		}
		tags = append(tags, existingTag)
	}

	links := []models.Link{}
	for _, linkURL := range markdownData.MetaData.Links {
		var existingLink models.Link
		if err := db.Where("url = ?", linkURL).First(&existingLink).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to query link: %w", err)
			}
			existingLink = models.Link{URL: linkURL}
			if err := db.Create(&existingLink).Error; err != nil {
				return nil, fmt.Errorf("failed to create link: %w", err)
			}
		}
		links = append(links, existingLink)
	}

	var langCode models.ArticleLangCode
	langCode, ok := langCode.FromString(markdownData.MetaData.Lang)
	if !ok {
		return nil, fmt.Errorf("unsupported language code in file %s", path)
	}

	createdAt, err := time.Parse(time.DateOnly, markdownData.MetaData.CreateAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CreateAt in file %s: %v", path, err)
	}

	return &models.Article{
		ID:           uint(markdownData.MetaData.ID),
		Title:        markdownData.MetaData.Title,
		Introduction: markdownData.MetaData.Introduction,
		Tags:         tags,
		Subject:      markdownData.Subject,
		CreatedAt:    createdAt,
		Lang:         langCode,
		Links:        links,
	}, nil
}
