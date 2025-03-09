package handler_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"share-basket-auth-service/presentation/api/middleware"
	"share-basket-auth-service/presentation/api/server"
)

func setupTestServer() *server.Server {
	s := server.New(":8080")
	s.Use(middleware.ErrorMiddleware())
	return s
}

func newJSONRequest(method, path string, body map[string]interface{}) *http.Request {
	reqBody, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	return req
}
