package comments

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mstyushin/go-news-api-gw/pkg/client"
	"github.com/mstyushin/go-news-comments/pkg/api"
	"github.com/mstyushin/go-news-comments/pkg/storage"
)

const (
	getComments = "/comments/by-articleid"
	addComment  = "/comments"
	moderate    = "/moderation"
)

type Client struct {
	CommentsSVC   string
	ModerationSVC string
	HttpClient    *client.HttpClient
}

func NewClient(commentsService, moderationService string) *Client {
	return &Client{
		CommentsSVC:   commentsService,
		ModerationSVC: moderationService,
		HttpClient:    client.New(),
	}
}

func (c *Client) AddComment(ctx context.Context, comment storage.Comment) (api.CommentCreatedResponse, error) {
	body, err := c.HttpClient.POST(ctx, fmt.Sprintf("http://%s%s", c.CommentsSVC, addComment), comment)
	if err != nil {
		log.Println(err.Error())
		return api.CommentCreatedResponse{}, err
	}

	var res api.CommentCreatedResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println(err.Error())
		return api.CommentCreatedResponse{}, err
	}

	return res, nil
}

func (c *Client) GetComments(ctx context.Context, articleID int) ([]storage.Comment, error) {
	body, err := c.HttpClient.GET(ctx, fmt.Sprintf("http://%s%s/%d", c.CommentsSVC, getComments, articleID))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var comments []storage.Comment
	err = json.Unmarshal(body, &comments)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return comments, nil
}

// Forward comment to moderation service.
// Return false either if moderation fails or any error occured in process.
func (c *Client) Moderate(ctx context.Context, comment storage.Comment) bool {
	_, err := c.HttpClient.POST(ctx, fmt.Sprintf("http://%s%s", c.ModerationSVC, moderate), comment)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}
