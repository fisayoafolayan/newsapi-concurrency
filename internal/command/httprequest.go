package command

import (
	"context"
	"fmt"
	"os"
	"time"

	"concurrency-demo/internal/entity"
	"concurrency-demo/internal/httpclient/news/apple"
	"concurrency-demo/internal/httpclient/news/techcrunch"
	"concurrency-demo/internal/httpclient/news/tesla"
	"concurrency-demo/pkg/bootstrap"
	"github.com/spf13/cobra"

	"golang.org/x/sync/errgroup"
)

func HTTPRequestsCommand() *cobra.Command {
	newsRetrieval := NewNewsRetrievalRequests()
	newsAPICommand := &cobra.Command{
		Use:   "news:retrieve",
		Short: "Retrieve news from news api",
		RunE: func(cmd *cobra.Command, args []string) error {
			return newsRetrieval.Run()
		},
	}

	return newsAPICommand
}

type NewsRetrievalRequest struct {
	appleNewsAPI      *apple.NewsAPI
	techcrunchNewsAPI *techcrunch.NewsAPI
	teslaNewsAPI      *tesla.NewsAPI
}

func NewNewsRetrievalRequests() *NewsRetrievalRequest {
	return &NewsRetrievalRequest{}
}

func (r *NewsRetrievalRequest) Run() error {
	return bootstrap.Container.Invoke(r.run)
}

func (r *NewsRetrievalRequest) run(
	AppleNewsAPI *apple.NewsAPI,
	TechcrunchNewsAPI *techcrunch.NewsAPI,
	TeslaNewsAPI *tesla.NewsAPI,
) {
	r.appleNewsAPI = AppleNewsAPI
	r.techcrunchNewsAPI = TechcrunchNewsAPI
	r.teslaNewsAPI = TeslaNewsAPI

	fmt.Println("Creating async requests...")
	r.asyncRequest()
	fmt.Println("Done.")

	fmt.Println("--------------------------")

	fmt.Println("Creating sync requests...")
	r.syncRequest()
	fmt.Println("Done.")

	os.Exit(0)
}

func (r *NewsRetrievalRequest) asyncRequest() {
	start := time.Now()

	var appleNewsResponse *entity.NewsResponse
	var techcrunchNewsResponse *entity.NewsResponse
	var teslaNewsResponse *entity.NewsResponse

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		var err error
		appleNewsResponse, err = r.appleNewsAPI.GetNews(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		techcrunchNewsResponse, err = r.techcrunchNewsAPI.GetNews(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		teslaNewsResponse, err = r.teslaNewsAPI.GetNews(ctx)
		return err
	})

	if err := g.Wait(); err != nil {
		fmt.Println("error", err)
		return
	}

	fmt.Println("--------------------------")

	fmt.Printf("--- Apple News API Status: %s\n", appleNewsResponse.Status)
	fmt.Printf("--- Apple News API Result: %v\n", appleNewsResponse.TotalResults)

	fmt.Println("--------------------------")

	fmt.Printf("--- Techcrunch News API Status: %s\n", techcrunchNewsResponse.Status)
	fmt.Printf("--- Techcrunch News API Result: %v\n", techcrunchNewsResponse.TotalResults)

	fmt.Println("--------------------------")

	fmt.Printf("--- Tesla News API Status: %s\n", teslaNewsResponse.Status)
	fmt.Printf("--- Tesla News API Result: %v\n", teslaNewsResponse.TotalResults)

	fmt.Println("--------------------------")

	end := time.Now()
	fmt.Printf("--- Completed after %v seconds\n", end.Sub(start).Seconds())
}

func (r *NewsRetrievalRequest) syncRequest() {
	start := time.Now()
	fmt.Println("--------------------------")

	appleNewsAPIResponse, err := r.appleNewsAPI.GetNews(context.Background())
	fmt.Printf("--- Apple News API Status: %s\n", appleNewsAPIResponse.Status)
	fmt.Printf("--- Apple News API Total Result: %v\n", appleNewsAPIResponse.TotalResults)

	fmt.Println("--------------------------")

	techcrunchNewsAPIResponse, err := r.techcrunchNewsAPI.GetNews(context.Background())
	fmt.Printf("--- Apple News API Status: %s\n", techcrunchNewsAPIResponse.Status)
	fmt.Printf("--- Apple News API Total Result: %v\n", techcrunchNewsAPIResponse.TotalResults)

	fmt.Println("--------------------------")

	teslaNewsAPIResponse, err := r.teslaNewsAPI.GetNews(context.Background())
	fmt.Printf("--- Tesla News API  Status: %s\n", teslaNewsAPIResponse.Status)
	fmt.Printf("--- Tesla News API Total Result: %v\n", teslaNewsAPIResponse.TotalResults)

	fmt.Println("--------------------------")

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	end := time.Now()
	fmt.Printf("--- Completed after %v seconds\n", end.Sub(start).Seconds())
}
