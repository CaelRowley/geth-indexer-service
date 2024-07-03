package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CaelRowley/geth-indexer-service/pkg/handlers"
	"github.com/stretchr/testify/assert"
)

func TestChiRouter(t *testing.T) {
	router := NewRouter(nil)
	handlers.Init(nil, nil, router)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	methods := []string{"POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	for _, method := range methods {
		req, err := http.NewRequest(method, "/", nil)
		assert.NoError(t, err)
		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.NotEqual(t, http.StatusOK, rr.Code)
	}
}
