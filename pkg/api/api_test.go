package api

import (
	"encoding/json"
	"fmt"
	"go-news-api-gw/pkg/config"
	"go-news-api-gw/pkg/model"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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
	req := httptest.NewRequest("GET", fmt.Sprintf("%s/news/latest?page=1", cfg.BaseURL), nil)
	w := httptest.NewRecorder()
	api.mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "cannot read response body")

	var articles model.ArticlesShort
	err = json.Unmarshal(body, &articles)
	assert.NoError(t, err, "cannot unmarshal response body")

	assert.Equal(t, fmt.Sprintf("%s/news/1", cfg.BaseURL), articles[0].LinkToFull)
}

func TestAPI_generateLinkToFull(t *testing.T) {
	a := model.ArticleShort{
		ID:           1,
		Title:        "Article 1",
		ShortContent: "Some meaningful content",
		PubTime:      1717882929,
	}
	api.generateLinkToFull(&a)

	assert.Equal(t, fmt.Sprintf("%s/news/%d", cfg.BaseURL, a.ID), a.LinkToFull)
}
