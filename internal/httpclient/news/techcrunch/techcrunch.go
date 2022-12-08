package techcrunch

import (
	"context"
	"net/http"
	"os"

	"newsapi-concurrency/internal/entity"
	"newsapi-concurrency/internal/httpclient"
	"newsapi-concurrency/internal/timeutil"
)

const (
	sortBy = "&sortBy=publishedAt"
)

type NewsAPI struct {
	httpClient *httpclient.HttpClient
}

func NewClient(httpClient *httpclient.HttpClient) *NewsAPI {
	return &NewsAPI{
		httpClient: httpClient,
	}
}

func (p *NewsAPI) GetNews(ctx context.Context) (*entity.NewsResponse, error) {
	var newsResponse *entity.NewsResponse
	var currentDate = timeutil.CurrentDate("2006-01-02")

	err := p.httpClient.Get(
		ctx,
		os.Getenv("NEWS_API_TECHCRUNCH_ENDPOINT")+
			"&from="+
			currentDate+sortBy,
		http.NoBody,
		&newsResponse,
	)
	if err != nil {
		return nil, err
	}
	return newsResponse, nil
}
