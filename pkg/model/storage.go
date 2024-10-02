// TODO import this from go-news-scraper and go-news-comments
package model

type ArticleShort struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	ShortContent string `json:"short_content"`
	LinkToFull   string `json:"link_to_full"`
	PubTime      int64  `json:"pub_time"`
}

type ArticleFull struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Link      string `json:"link"`
	RSSFeedID int    `json:"rss_feed_id"`
	PubTime   int64  `json:"pub_time"`
}

type Comment struct {
	ID        int    `json:"id"`
	ArticleID int    `json:"article_id"`
	ParentID  int    `json:"parent_id"`
	Author    string `json:"author"`
	Text      string `json:"text"`
	PubTime   int64  `json:"pub_time"`
}
