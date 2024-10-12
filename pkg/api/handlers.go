package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	storageComments "github.com/mstyushin/go-news-comments/pkg/storage"
	scraperAPI "github.com/mstyushin/go-news-scraper/pkg/api"
	storageScraper "github.com/mstyushin/go-news-scraper/pkg/storage"

	"github.com/gorilla/mux"
)

// TODO make it configurable
var pageSize = 10

func (api *API) getNews(w http.ResponseWriter, r *http.Request) {
	var s string
	pageNum := 1

	if r.URL.Query().Has("page_size") {
		s = r.URL.Query().Get("page_size")
		pageSize, _ = strconv.Atoi(s)
	}

	if r.URL.Query().Has("page") {
		s = r.URL.Query().Get("page")
		pageNum, _ = strconv.Atoi(s)
	}

	var paginated scraperAPI.PaginatedResponse
	var err error

	if r.URL.Query().Has("s") {
		searchString := r.URL.Query().Get("s")
		paginated, err = api.scraperClient.SearchNews(r.Context(), searchString, pageSize, pageNum)
	} else {
		paginated, err = api.scraperClient.GetNews(r.Context(), pageSize, pageNum)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for idx, _ := range paginated.Articles {
		api.generateLinkToFull(&paginated.Articles[idx])
	}

	bytes, err := json.Marshal(paginated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (api *API) getDetailedArticle(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(s)

	// TODO magic number
	results := make(chan interface{}, 2)

	// two goroutines for 2 services
	// TODO magic number
	wg.Add(2)

	// get news from Scraper service
	go func() {
		log.Println("Getting data from Scraper service")
		defer wg.Done()
		article, err := api.scraperClient.GetArticle(r.Context(), id)
		if err != nil {
			results <- err
			return
		}
		results <- article
	}()

	// get comments from Comments service
	go func() {
		log.Println("Getting data from Comments service")
		defer wg.Done()
		comments, err := api.commentsClient.GetComments(r.Context(), id)
		if err != nil {
			results <- err
			return
		}
		results <- comments
	}()

	wg.Wait()
	close(results)

	response := ArticleFullResponse{}

	for data := range results {
		switch data.(type) {
		case error:
			log.Println("Got error from one of the services")
			http.Error(w, data.(error).Error(), http.StatusInternalServerError)
		case storageScraper.Article:
			response.Article = data.(storageScraper.Article)
		case []storageComments.Comment:
			response.Comments = data.([]storageComments.Comment)
		}
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (api *API) addComment(w http.ResponseWriter, r *http.Request) {
	shouldAddComment := r.URL.Query().Get(commentQueryParam)
	check, err := strconv.ParseBool(shouldAddComment)
	if err != nil {
		http.Error(w, fmt.Sprintf("expecting ?%s=true query param", commentQueryParam), http.StatusBadRequest)
		return
	}
	if !check {
		http.Error(w, fmt.Sprintf("query parameter: ?%s= must be true for POST requests", commentQueryParam), http.StatusNotAcceptable)
		return
	}

	var comment storageComments.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !api.commentsClient.Moderate(r.Context(), comment) {
		log.Println("Comment not allowed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s := mux.Vars(r)["id"]
	aid, _ := strconv.Atoi(s)
	comment.ArticleID = aid
	// TODO check if Article with ID=aid exists
	res, err := api.commentsClient.AddComment(r.Context(), comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}
