package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
)

const contentTypeJSON = "application/json"

type errorResponse struct {
	Errors struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

var ErrNotFound = errors.New("not found")

type Headers map[string]string

type HttpClient struct {
	http.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}

func (c *HttpClient) Get(ctx context.Context, endpoint string, body interface{}, responseEntity interface{}) error {
	return c.Request(ctx, http.MethodGet, endpoint, body, responseEntity)
}

func (c *HttpClient) Request(ctx context.Context, method, endpoint string, body interface{}, responseEntity interface{}) error {

	b, err := json.Marshal(body)
	if err != nil {
		return errors.Wrapf(err, "prepare Request body for %s %s", method, c.generateEndpointPath(endpoint))
	}

	headers := Headers{
		"Content-Type": contentTypeJSON,
		"Accept":       contentTypeJSON,
	}

	response, err := c.send(ctx, method, endpoint, bytes.NewReader(b), headers)

	if err != nil {
		return errors.WithStack(err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if !isResponseSuccessful(response) {
		return buildResponseError(response)
	}

	if response.StatusCode == http.StatusNoContent {
		return nil
	}

	if responseEntity == nil {
		return nil
	}

	err = json.NewDecoder(response.Body).Decode(&responseEntity)
	if err != nil {
		return errors.Wrapf(err, "decode json response from %s %s", method, endpoint)
	}

	return nil
}

func (c *HttpClient) send(ctx context.Context, method, endpoint string, body io.Reader, headers Headers) (*http.Response, error) {
	url := c.generateEndpointPath(endpoint)
	request, err := http.NewRequestWithContext(ctx, method, url, body)

	if err != nil {
		return nil, errors.Wrapf(err, "build Request for %s", url)
	}

	for name, value := range headers {
		request.Header.Set(name, value)
	}

	response, err := c.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "Request %s %s", method, url)
	}

	return response, nil
}

func (c *HttpClient) generateEndpointPath(endpoint string) string {
	return fmt.Sprintf("%s%s&apiKey=%s", os.Getenv("NEWS_API_BASE_URI"), endpoint, os.Getenv("NEWS_API_APP_KEY"))
}

func isResponseSuccessful(response *http.Response) bool {
	return response.StatusCode < 400
}

func buildResponseError(response *http.Response) error {
	var errResponse *errorResponse
	var errData string

	contentType := strings.ToLower(response.Header.Get("Content-Type"))
	if contentType == contentTypeJSON {
		err := json.NewDecoder(response.Body).Decode(&errResponse)
		if err == nil {
			errData = fmt.Sprintf(" (%d: %s)", errResponse.Errors.Code, errResponse.Errors.Message)
		}
	}

	err := errors.Errorf("Request %s %s: got response with status %d%s",
		response.Request.Method,
		response.Request.URL.String(),
		response.StatusCode,
		errData,
	)

	if response.StatusCode == http.StatusNotFound {
		return errors.Wrap(ErrNotFound, err.Error())
	}

	return err
}
