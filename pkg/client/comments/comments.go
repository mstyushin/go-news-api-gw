package comments

import (
	"context"
	"encoding/json"
	"fmt"
	"go-news-api-gw/pkg/client"
	"go-news-api-gw/pkg/model"
	"log"
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

func (c *Client) AddComment(ctx context.Context, comment model.Comment) (model.CommentCreatedResponse, error) {
	body, err := c.HttpClient.POST(ctx, fmt.Sprintf("http://%s%s", c.CommentsSVC, addComment), comment)
	if err != nil {
		log.Println(err.Error())
		return model.CommentCreatedResponse{}, err
	}

	var res model.CommentCreatedResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println(err.Error())
		return model.CommentCreatedResponse{}, err
	}

	return res, nil
}

func (c *Client) GetComments(ctx context.Context, articleID int) ([]model.Comment, error) {
	body, err := c.HttpClient.GET(ctx, fmt.Sprintf("http://%s%s/%d", c.CommentsSVC, getComments, articleID))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var comments []model.Comment
	err = json.Unmarshal(body, &comments)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return comments, nil
}

// Forward comment to moderation service.
// Return false either if moderation fails or any error occured in process.
func (c *Client) Moderate(ctx context.Context, comment model.Comment) bool {
	_, err := c.HttpClient.POST(ctx, fmt.Sprintf("http://%s%s", c.ModerationSVC, moderate), comment)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}
