package external

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

type APIClient interface {
	MakeAPIRequest(ctx context.Context, url string, queryParams map[string]string, headers map[string]string) (*http.Response, error)
}

type ApiClientBase struct {
	Logger *logrus.Logger
}

func NewAPIClient(logger *logrus.Logger) *ApiClientBase {
	return &ApiClientBase{
		Logger: logger,
	}
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

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	c.Logger.WithFields(logrus.Fields{
		"requestURL":     request.URL.String(),
		"requestMethod":  request.Method,
		"requestHeaders": request.Header,
	}).Info("Making API request")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	c.Logger.WithFields(logrus.Fields{
		"responseStatusCode": resp.StatusCode,
		"responseHeaders":    resp.Header,
	}).Info("Received API response")

	return resp, nil
}
