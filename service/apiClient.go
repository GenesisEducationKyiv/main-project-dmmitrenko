package service

import (
	"context"
	"net/http"
)

type APIClient interface {
	MakeAPIRequest(ctx context.Context, url string, queryParams map[string]string, headers map[string]string) (*http.Response, error)
}

type ApiClientBase struct{}

func NewAPIClient() ApiClientBase {
	return ApiClientBase{}
}

func (c *ApiClientBase) MakeAPIRequest(ctx context.Context, url string, queryParams map[string]string, headers map[string]string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)

	q := request.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	request.URL.RawQuery = q.Encode()

	if headers != nil {
		for key, value := range headers {
			request.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
