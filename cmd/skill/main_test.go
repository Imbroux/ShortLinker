package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_webhook(t *testing.T) {

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
			url:          "/",
			headers:      map[string]string{"Content-Type": "text/plain"},
			body:         "https://practicum.yandex.ru/",
			expectedCode: http.StatusCreated,
			expectedBody: "",
		},
		{
			method:       http.MethodGet,
			url:          "http://localhost:8080",
			headers:      nil,
			body:         "",
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.method, func(t *testing.T) {
			bodyReader := strings.NewReader(tc.body)
			r := httptest.NewRequest(tc.method, tc.url, bodyReader)
			for key, value := range tc.headers {
				r.Header.Set(key, value)
			}
			w := httptest.NewRecorder()

			webhook(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}

		})
	}
}
