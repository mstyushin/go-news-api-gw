package model

type CommentCreatedResponse struct {
	ID      int   `json:"id"`
	PubTime int64 `json:"pub_time"`
}

type ArticleFullResponse struct {
	Article  ArticleFull `json:"article"`
	Comments []Comment   `json:"comments"`
}

type Paginator struct {
	NumPages int `json:"num_pages"`
	CurPage  int `json:"cur_page"`
	PageSize int `json:"page_size"`
}

type PaginatedResponse struct {
	Articles  []ArticleShort `json:"articles"`
	Paginator Paginator      `json:"paginator"`
}
