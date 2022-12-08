package provider

import (
	"newsapi-concurrency/internal/httpclient"
	"newsapi-concurrency/internal/httpclient/news/apple"
	"newsapi-concurrency/internal/httpclient/news/techcrunch"
	"newsapi-concurrency/internal/httpclient/news/tesla"
	"newsapi-concurrency/pkg/bootstrap"

	"github.com/labstack/echo"
)

type Provider struct {
	Echo *echo.Echo
}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Provide() {
	bootstrap.Provide(httpclient.NewHttpClient)
	bootstrap.Provide(apple.NewClient)
	bootstrap.Provide(techcrunch.NewClient)
	bootstrap.Provide(tesla.NewClient)
}
