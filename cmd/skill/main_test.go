package main

import (
	"YandexLearnMiddle/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_webhook(t *testing.T) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", handlers.HandlePost)
	r.Get("/*", handlers.HandleGet())

	srv := httptest.NewServer(r)
	defer srv.Close()

	tests := []struct {
		method       string
		url          string
		headers      map[string]string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			method:       http.MethodPost,
			url:          srv.URL + "/api/shorten",
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         `{"url":"https://practicum.yandex.ru/"}`,
			expectedCode: http.StatusCreated,
			expectedBody: `{"result":"http://localhost:8080/`,
		},
		{
			method:       http.MethodGet,
			url:          srv.URL + "/nonexistent",
			headers:      nil,
			body:         "",
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
	}

	client := resty.New()

	for _, tc := range tests {
		t.Run(tc.method, func(t *testing.T) {
			req := client.R().
				SetHeaders(tc.headers).
				SetBody(tc.body)

			resp, err := req.Execute(tc.method, tc.url)

			assert.NoError(t, err, "Unexpected error during request")

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code does not match expected")

			if tc.expectedBody != "" {
				assert.Contains(t, resp.String(), tc.expectedBody, "Response body does not contain expected result")
			}
		})
	}
}
