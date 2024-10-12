package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mstyushin/go-news-api-gw/pkg/config"
	scraperAPI "github.com/mstyushin/go-news-scraper/pkg/api"
	"github.com/mstyushin/go-news-scraper/pkg/storage"

	"github.com/stretchr/testify/assert"
)

var api *API
var cfg *config.Config

func TestMain(m *testing.M) {
	cfg = config.DefaultConfig()
	api = New(cfg)

	os.Exit(m.Run())
}

func TestAPI_getNews(t *testing.T) {
	req := httptest.NewRequest("GET", fmt.Sprintf("%s/news/latest", cfg.BaseURL), nil)
	rr := httptest.NewRecorder()

	api.mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, rr.Header().Get("x-request-id") != "", "should populate x-request-id header")

	b, err := ioutil.ReadAll(rr.Body)
	assert.NoError(t, err, "should be able to read response body")

	var articles scraperAPI.PaginatedResponse
	err = json.Unmarshal(b, &articles)
	assert.NoError(t, err, "cannot unmarshal response body")
	assert.Equal(t, fmt.Sprintf("%s/news/1", cfg.BaseURL), articles.Articles[0].LinkToFull, "expecting correct link to full article")
}

func TestAPI_generateLinkToFull(t *testing.T) {
	a := storage.ArticleShort{
		ID:           1,
		Title:        "Article 1",
		ShortContent: "Some meaningful content",
		PubTime:      1717882929,
	}
	api.generateLinkToFull(&a)

	assert.Equal(t, fmt.Sprintf("%s/news/%d", cfg.BaseURL, a.ID), a.LinkToFull)
}
