package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpService interface {
	Post(ctx context.Context, path string, body []byte, headers map[string]any) (*http.Response, error)
	Delete(ctx context.Context, path string, body []byte, headers map[string]any) (*http.Response, error)
}

type BaseService struct {
	Client  *http.Client
	BaseURL string
	Headers map[string]string
}

func NewBaseService(baseURL string) *BaseService {
	return &BaseService{
		Client: &http.Client{
			Timeout: 2 * time.Minute,
		},
		BaseURL: baseURL,
	}
}

func (s *BaseService) Post(ctx context.Context, path string, body []byte) (*http.Response, int, error) {

	url := s.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range s.Headers {
		req.Header.Set(key, fmt.Sprintf("%v", value))
	}
	response, err := s.Client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer response.Body.Close()

	bodyResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, response.StatusCode, fmt.Errorf("status code: %d, response: %s", response.StatusCode, string(bodyResponse))
	}

	// Create a new response with the body
	newResponse := &http.Response{
		StatusCode: response.StatusCode,
		Body:       io.NopCloser(bytes.NewBuffer(bodyResponse)),
		Header:     response.Header,
	}

	return newResponse, http.StatusOK, nil
}

func (s *BaseService) Delete(ctx context.Context, path string, body []byte) (*http.Response, int, error) {

	url := s.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range s.Headers {
		req.Header.Set(key, fmt.Sprintf("%v", value))
	}
	response, err := s.Client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return response, http.StatusOK, nil
}
