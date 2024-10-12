package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/mstyushin/go-news-api-gw/pkg/client/comments"
	"github.com/mstyushin/go-news-api-gw/pkg/client/scraper"
	"github.com/mstyushin/go-news-api-gw/pkg/config"
	"github.com/mstyushin/go-news-scraper/pkg/storage"

	"github.com/gorilla/mux"
)

const (
	pageQueryParam     = "page"
	pageSizeQueryParam = "page_size"
	commentQueryParam  = "c"
)

var wg sync.WaitGroup

type API struct {
	HttpListenPort int
	mux            *mux.Router
	commentsClient *comments.Client
	scraperClient  *scraper.Client
	baseURL        string
}

func New(cfg *config.Config) *API {
	api := API{
		HttpListenPort: cfg.HttpPort,
		mux:            mux.NewRouter(),
		baseURL:        cfg.BaseURL,
		commentsClient: comments.NewClient(
			fmt.Sprintf("%s:%d", cfg.CommentsServiceAddress, cfg.CommentsServicePort),
			fmt.Sprintf("%s:%d", cfg.ModerationServiceAddress, cfg.ModerationServicePort),
		),
		scraperClient: scraper.NewClient(fmt.Sprintf("%s:%d", cfg.NewsScraperAddress, cfg.NewsScraperPort)),
	}

	api.endpoints()
	return &api
}

func (api *API) Run(ctx context.Context) error {
	errChan := make(chan error)
	srv := api.serve(ctx, errChan)

	select {
	case <-ctx.Done():
		log.Println("gracefully shutting down")
		srv.Shutdown(ctx)
		return ctx.Err()
	case err := <-errChan:
		log.Println(err)
		return err
	}
}

func (api *API) serve(ctx context.Context, errChan chan error) *http.Server {
	log.Println("serving HTTP server at", api.HttpListenPort)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", api.HttpListenPort),
		Handler: api.mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, fmt.Sprintf(":%v", api.HttpListenPort), l.Addr().String())
			return ctx
		},
	}

	go func(s *http.Server) {
		if err := s.ListenAndServe(); err != nil {
			errChan <- err
		}
	}(httpServer)

	return httpServer
}

func (api *API) endpoints() {
	api.mux.HandleFunc("/news/latest", api.getNews).Methods(http.MethodGet, http.MethodOptions)
	api.mux.HandleFunc("/news/{id}", api.getDetailedArticle).Methods(http.MethodGet, http.MethodOptions)
	api.mux.HandleFunc("/news/{id}", api.addComment).Methods(http.MethodPost, http.MethodOptions)
	api.mux.Use(URLSchemaMiddleware(api.mux))
	api.mux.Use(RequestIDLoggerMiddleware(api.mux))
	api.mux.Use(LoggerMiddleware(api.mux))
}

func (api *API) generateLinkToFull(a *storage.ArticleShort) {
	a.LinkToFull = fmt.Sprintf("%s/news/%d", api.baseURL, a.ID)
}
