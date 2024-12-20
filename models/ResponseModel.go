package models

type ArticleDTO struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Introduction string    `json:"introduction"`
	CreateAt     string    `json:"createAt"`
	Subject      string    `json:"subject"`
	Lang         int       `json:"lang"`
	Tags         []TagDTO  `json:"tags"`
	Links        []LinkDTO `json:"links"`
}

type TagDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type LinkDTO struct {
	URL string `json:"url"`
}
