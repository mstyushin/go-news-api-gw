package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mstyushin/go-news-api-gw/pkg/client"
	"github.com/mstyushin/go-news-scraper/pkg/api"
	"github.com/mstyushin/go-news-scraper/pkg/storage"
)

const (
	getNews = "/news"
)

type Client struct {
	ScraperSVC string
	HttpClient *client.HttpClient
}

func NewClient(scraperService string) *Client {
	return &Client{
		ScraperSVC: scraperService,
		HttpClient: client.New(),
	}
}

func (c *Client) GetNews(ctx context.Context, pageSize, pageNum int) (api.PaginatedResponse, error) {
	url := fmt.Sprintf("http://%s%s?page_size=%d&page=%d", c.ScraperSVC, getNews, pageSize, pageNum)
	body, err := c.HttpClient.GET(ctx, url)
	if err != nil {
		log.Println(err.Error())
		return api.PaginatedResponse{}, err
	}

	var paginated api.PaginatedResponse
	err = json.Unmarshal(body, &paginated)
	if err != nil {
		log.Println(err.Error())
		return api.PaginatedResponse{}, err
	}

	return paginated, nil
}

func (c *Client) SearchNews(ctx context.Context, searchString string, pageSize, pageNum int) (api.PaginatedResponse, error) {
	// TODO think about DRYing this code
	url := fmt.Sprintf("http://%s%s?s=%s&page_size=%d&page=%d", c.ScraperSVC, getNews, searchString, pageSize, pageNum)
	body, err := c.HttpClient.GET(ctx, url)
	if err != nil {
		log.Println(err.Error())
		return api.PaginatedResponse{}, err
	}

	var paginated api.PaginatedResponse
	err = json.Unmarshal(body, &paginated)
	if err != nil {
		log.Println(err.Error())
		return api.PaginatedResponse{}, err
	}

	return paginated, nil
}

func (c *Client) GetArticle(ctx context.Context, id int) (storage.Article, error) {
	url := fmt.Sprintf("http://%s%s/%d", c.ScraperSVC, getNews, id)
	body, err := c.HttpClient.GET(ctx, url)
	if err != nil {
		log.Println(err.Error())
		return storage.Article{}, err
	}

	var article storage.Article
	err = json.Unmarshal(body, &article)
	if err != nil {
		log.Println(err.Error())
		return storage.Article{}, err
	}

	return article, nil
}
