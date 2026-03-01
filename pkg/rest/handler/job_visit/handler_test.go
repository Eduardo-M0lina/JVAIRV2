package job_visit

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	domainFile "github.com/your-org/jvairv2/pkg/domain/file"
	domainJV "github.com/your-org/jvairv2/pkg/domain/job_visit"
)

// Mock implementations
type mockJobVisitUseCase struct {
	visits []*domainJV.JobVisit
}

func (m *mockJobVisitUseCase) Create(ctx context.Context, visit *domainJV.JobVisit) (int64, error) {
	visit.ID = int64(len(m.visits) + 1)
	m.visits = append(m.visits, visit)
	return visit.ID, nil
}

func (m *mockJobVisitUseCase) GetByID(ctx context.Context, id int64) (*domainJV.JobVisit, error) {
	for _, v := range m.visits {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, domainJV.ErrJobVisitNotFound
}

func (m *mockJobVisitUseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int, sort, direction string) ([]*domainJV.JobVisit, int64, error) {
	return m.visits, int64(len(m.visits)), nil
}

func (m *mockJobVisitUseCase) Update(ctx context.Context, visit *domainJV.JobVisit) error {
	for i, v := range m.visits {
		if v.ID == visit.ID {
			m.visits[i] = visit
			return nil
		}
	}
	return domainJV.ErrJobVisitNotFound
}

func (m *mockJobVisitUseCase) Delete(ctx context.Context, id int64) error {
	for i, v := range m.visits {
		if v.ID == id {
			m.visits = append(m.visits[:i], m.visits[i+1:]...)
			return nil
		}
	}
	return domainJV.ErrJobVisitNotFound
}

type mockFileUseCase struct {
	files []*domainFile.File
}

func (m *mockFileUseCase) Upload(ctx context.Context, fileableID int64, fileableType string, filename string, contentType string, body io.Reader) (*domainFile.File, error) {
	f := &domainFile.File{
		ID:           int64(len(m.files) + 1),
		Type:         &contentType,
		URL:          "https://example.com/" + filename,
		FileableID:   fileableID,
		FileableType: fileableType,
	}
	m.files = append(m.files, f)
	return f, nil
}

func (m *mockFileUseCase) GetByID(ctx context.Context, id int64) (*domainFile.File, error) {
	for _, f := range m.files {
		if f.ID == id {
			return f, nil
		}
	}
	return nil, domainFile.ErrFileNotFound
}

func (m *mockFileUseCase) ListByFileable(ctx context.Context, fileableID int64, fileableType string) ([]*domainFile.File, error) {
	var result []*domainFile.File
	for _, f := range m.files {
		if f.FileableID == fileableID && f.FileableType == fileableType {
			result = append(result, f)
		}
	}
	return result, nil
}

func (m *mockFileUseCase) Delete(ctx context.Context, id int64) error {
	for i, f := range m.files {
		if f.ID == id {
			m.files = append(m.files[:i], m.files[i+1:]...)
			return nil
		}
	}
	return domainFile.ErrFileNotFound
}

func (m *mockFileUseCase) Download(ctx context.Context, id int64) (io.ReadCloser, string, string, error) {
	for _, f := range m.files {
		if f.ID == id {
			return io.NopCloser(strings.NewReader("test content")), "text/plain", "test.txt", nil
		}
	}
	return nil, "", "", domainFile.ErrFileNotFound
}

func setupHandler() *Handler {
	jobVisitUC := &mockJobVisitUseCase{}
	fileUC := &mockFileUseCase{}
	return NewHandler(jobVisitUC, fileUC)
}

func TestHandler_List(t *testing.T) {
	h := setupHandler()

	// Set up chi context with jobId
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	req := httptest.NewRequest("GET", "/jobs/1/visits?page=1&pageSize=10", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.List(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandler_Get_Success(t *testing.T) {
	h := setupHandler()

	// Create a visit first
	ctx := context.Background()
	visit := &domainJV.JobVisit{
		JobID:  1,
		UserID: 1,
		Date:   time.Now(),
	}
	_, _ = h.useCase.Create(ctx, visit)

	// Set up chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	rctx.URLParams.Add("visitId", "1")
	req := httptest.NewRequest("GET", "/jobs/1/visits/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Get(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandler_Get_NotFound(t *testing.T) {
	h := setupHandler()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	rctx.URLParams.Add("visitId", "999")
	req := httptest.NewRequest("GET", "/jobs/1/visits/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Get(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestHandler_Create_Success(t *testing.T) {
	h := setupHandler()

	body := `{
		"userId": 1,
		"viewableBy": "[\"1\",\"2\"]",
		"date": "01-15-2024",
		"report": "Test visit report"
	}`

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	req := httptest.NewRequest("POST", "/jobs/1/visits", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
}

func TestHandler_Create_InvalidBody(t *testing.T) {
	h := setupHandler()

	body := `invalid json`

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	req := httptest.NewRequest("POST", "/jobs/1/visits", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestHandler_Delete_Success(t *testing.T) {
	h := setupHandler()

	// Create a visit first
	ctx := context.Background()
	visit := &domainJV.JobVisit{
		JobID:  1,
		UserID: 1,
		Date:   time.Now(),
	}
	_, _ = h.useCase.Create(ctx, visit)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	rctx.URLParams.Add("visitId", "1")
	req := httptest.NewRequest("DELETE", "/jobs/1/visits/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}

func TestHandler_Delete_NotFound(t *testing.T) {
	h := setupHandler()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	rctx.URLParams.Add("visitId", "999")
	req := httptest.NewRequest("DELETE", "/jobs/1/visits/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestHandler_ListFiles(t *testing.T) {
	h := setupHandler()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	rctx.URLParams.Add("visitId", "1")
	req := httptest.NewRequest("GET", "/jobs/1/visits/1/files", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.ListFiles(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandler_UploadFile_Success(t *testing.T) {
	h := setupHandler()

	// Create a visit first
	ctx := context.Background()
	visit := &domainJV.JobVisit{
		JobID:  1,
		UserID: 1,
		Date:   time.Now(),
	}
	_, _ = h.useCase.Create(ctx, visit)

	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("files", "test.txt")
	_, _ = part.Write([]byte("test content"))
	_ = writer.Close()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	rctx.URLParams.Add("visitId", "1")
	req := httptest.NewRequest("POST", "/jobs/1/visits/1/files", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.UploadFile(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
}

func TestHandler_UploadFile_NoFiles(t *testing.T) {
	h := setupHandler()

	// Create a visit first
	ctx := context.Background()
	visit := &domainJV.JobVisit{
		JobID:  1,
		UserID: 1,
		Date:   time.Now(),
	}
	_, _ = h.useCase.Create(ctx, visit)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	rctx.URLParams.Add("visitId", "1")
	req := httptest.NewRequest("POST", "/jobs/1/visits/1/files", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.UploadFile(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestHandler_DownloadFile_Success(t *testing.T) {
	h := setupHandler()

	// Create a file first
	ctx := context.Background()
	_, _ = h.fileUseCase.Upload(ctx, 1, fileableType, "test.txt", "text/plain", strings.NewReader("test"))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	rctx.URLParams.Add("visitId", "1")
	rctx.URLParams.Add("fileId", "1")
	req := httptest.NewRequest("GET", "/jobs/1/visits/1/files/1/download", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.DownloadFile(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "text/plain" {
		t.Errorf("Expected Content-Type text/plain, got %s", resp.Header.Get("Content-Type"))
	}
}

func TestHandler_DownloadFile_NotFound(t *testing.T) {
	h := setupHandler()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	rctx.URLParams.Add("visitId", "1")
	rctx.URLParams.Add("fileId", "999")
	req := httptest.NewRequest("GET", "/jobs/1/visits/1/files/999/download", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.DownloadFile(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestHandler_DeleteFile_Success(t *testing.T) {
	h := setupHandler()

	// Create a file first
	ctx := context.Background()
	_, _ = h.fileUseCase.Upload(ctx, 1, fileableType, "test.txt", "text/plain", strings.NewReader("test"))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	rctx.URLParams.Add("visitId", "1")
	rctx.URLParams.Add("fileId", "1")
	req := httptest.NewRequest("DELETE", "/jobs/1/visits/1/files/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.DeleteFile(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}

func TestHandler_DeleteFile_NotFound(t *testing.T) {
	h := setupHandler()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "1")
	rctx.URLParams.Add("visitId", "1")
	rctx.URLParams.Add("fileId", "999")
	req := httptest.NewRequest("DELETE", "/jobs/1/visits/1/files/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.DeleteFile(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}
