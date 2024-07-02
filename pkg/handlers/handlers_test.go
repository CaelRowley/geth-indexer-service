package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var handlers = Handlers{
	dbConn:    nil,
	ethClient: nil,
}

func TestMakeHandlerAPIErr(t *testing.T) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("test error"))
	}
	handler := makeHandler(handlerFunc)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := `{"statusCode":400,"msg":"test error"}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestMakeHandlerErr(t *testing.T) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		return fmt.Errorf("test error")
	}
	handler := makeHandler(handlerFunc)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := `{"statusCode":500,"msg":"interal server error"}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestHealthCheckHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handlers.healthCheckHandler))
	resp, err := http.Get(server.URL)
	assert.NoError(t, err)
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Healthy!", string(b))
}

func TestHealthCheckHandlerRR(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	handlers.healthCheckHandler(rr, req)

	defer rr.Result().Body.Close()
	b, err := io.ReadAll(rr.Result().Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Equal(t, "Healthy!", string(b))
}

func TestSetJSONResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"message": "test"}

	err := setJSONResponse(rr, http.StatusOK, data)
	assert.NoError(t, err)

	expected, _ := json.Marshal(data)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, string(expected), string(rr.Body.Bytes()))
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestAPIError(t *testing.T) {
	err := NewAPIError(http.StatusNotFound, fmt.Errorf("not found"))
	assert.Equal(t, "api error: 404", err.Error())
	assert.Equal(t, http.StatusNotFound, err.StatusCode)
	assert.Equal(t, "not found", err.Msg)

	invalidJSONErr := InvalidJson()
	assert.Equal(t, "invalid JSON request data", invalidJSONErr.Msg)
	assert.Equal(t, http.StatusBadRequest, invalidJSONErr.StatusCode)

	invalidURLParamErr := InvalidURLParam()
	assert.Equal(t, "invalid URLParam", invalidURLParamErr.Msg)
	assert.Equal(t, http.StatusBadRequest, invalidURLParamErr.StatusCode)

	invalidRequestDataErr := InvalidRequestData(map[string]string{"field": "error"})
	assert.Equal(t, map[string]string{"field": "error"}, invalidRequestDataErr.Msg)
	assert.Equal(t, http.StatusUnprocessableEntity, invalidRequestDataErr.StatusCode)
}
