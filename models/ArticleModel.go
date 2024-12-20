package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type ArticleLangCode int

const (
	ZH ArticleLangCode = 0 // 中文
	EN ArticleLangCode = 1 // 英文
	JP ArticleLangCode = 2 // 日文
)

// 实现 Stringer 接口（可选）
func (lang ArticleLangCode) String() string {
	switch lang {
	case ZH:
		return "ZH"
	case EN:
		return "EN"
	case JP:
		return "JP"
	default:
		return "UNKNOWN"
	}
}

func (lang ArticleLangCode) FromString(s string) (ArticleLangCode, bool) {
	s = strings.ToUpper(s)
	switch s {
	case "ZH":
		return ZH, true
	case "EN":
		return EN, true
	case "JP":
		return JP, true
	default:
		return 0, false
	}
}

func (lang ArticleLangCode) GetStringIterators() func() (string, bool) {
	langs := []ArticleLangCode{ZH, EN, JP}
	index := 0
	return func() (string, bool) {
		if index >= len(langs) {
			return "", false
		}
		val := langs[index].String()
		index++
		return val, true
	}
}

func (lang ArticleLangCode) GetValueIterators() func() (int, bool) {
	langs := []ArticleLangCode{ZH, EN, JP}
	index := 0
	return func() (int, bool) {
		if index >= len(langs) {
			return -1, false
		}
		val := langs[index]
		index++
		return int(val), true
	}
}

type Article struct {
	gorm.Model
	ID           uint `gorm:"primaryKey"`
	Title        string
	Introduction string
	Tags         []Tag `gorm:"many2many:article_tags;"`
	Subject      string
	CreatedAt    time.Time
	Lang         ArticleLangCode
	Links        []Link `gorm:"many2many:article_links;"`
}

type ArticleCountOfLang struct {
	Lang         string
	ArticleCount int64
}

type FrontMatter struct {
	ID           int      `yaml:"id"`
	Title        string   `yaml:"title"`
	Introduction string   `yaml:"introduction"`
	Tags         []string `yaml:"tags"`
	CreateAt     string   `yaml:"create_at"`
	Lang         string   `yaml:"lang"`
	Links        []string `yaml:"links"`
}

type ArticleMarkdownData struct {
	MetaData FrontMatter
	Subject  string
}
