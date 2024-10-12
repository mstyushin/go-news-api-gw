package api

import (
	storageComments "github.com/mstyushin/go-news-comments/pkg/storage"
	storageScraper "github.com/mstyushin/go-news-scraper/pkg/storage"
)

type ArticleFullResponse struct {
	Article  storageScraper.Article    `json:"article"`
	Comments []storageComments.Comment `json:"comments"`
}
