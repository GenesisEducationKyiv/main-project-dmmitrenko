package service

import "net/http"

type APIClient struct{}

func NewAPIClient() *APIClient {
	return &APIClient{}
}

func (c *APIClient) MakeAPIRequest(url string, queryParams map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
