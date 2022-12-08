package provider

import (
	"concurrency-demo/internal/httpclient"
	"concurrency-demo/internal/httpclient/news/apple"
	"concurrency-demo/internal/httpclient/news/techcrunch"
	"concurrency-demo/internal/httpclient/news/tesla"
	"concurrency-demo/pkg/bootstrap"

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
