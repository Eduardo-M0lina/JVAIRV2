package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockDBChecker struct {
	err error
}

func (m *mockDBChecker) Ping() error {
	return m.err
}

func TestHealthHandler_Check_Success(t *testing.T) {
	h := NewHealthHandler(&mockDBChecker{})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	h.Check(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp HealthResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
	assert.Equal(t, "ok", resp.DBStatus)
}

func TestHealthHandler_Check_DBError(t *testing.T) {
	h := NewHealthHandler(&mockDBChecker{err: assert.AnError})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	h.Check(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp HealthResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
}

func TestHealthHandler_Check_NilDBChecker(t *testing.T) {
	h := NewHealthHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	h.Check(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp HealthResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
}
