package service

import (
	"context"
	"net/http"
)

type APIClient struct{}

func NewAPIClient() APIClient {
	return APIClient{}
}

func (c *APIClient) MakeAPIRequest(ctx context.Context, url string, queryParams map[string]string) (*http.Response, error) {
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

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
