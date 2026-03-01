package job_visit

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	domainFile "github.com/your-org/jvairv2/pkg/domain/file"
	domainJV "github.com/your-org/jvairv2/pkg/domain/job_visit"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Handler maneja las peticiones HTTP para visitas de trabajo
type Handler struct {
	useCase     domainJV.Service
	fileUseCase domainFile.Service
}

// NewHandler crea una nueva instancia del handler de visitas de trabajo
func NewHandler(useCase domainJV.Service, fileUseCase domainFile.Service) *Handler {
	return &Handler{
		useCase:     useCase,
		fileUseCase: fileUseCase,
	}
}

// RegisterRoutes registra las rutas del handler
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/jobs/{jobId}/visits", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{visitId}", h.Get)
		r.Put("/{visitId}", h.Update)
		r.Delete("/{visitId}", h.Delete)
		// Archivos de visita
		r.Get("/{visitId}/files", h.ListFiles)
		r.Post("/{visitId}/files", h.UploadFile)
		r.Delete("/{visitId}/files/{fileId}", h.DeleteFile)
		r.Get("/{visitId}/files/{fileId}/download", h.DownloadFile)
	})
}

const timeFormat = "2006-01-02T15:04:05Z07:00"
const dateFormat = "01-02-2006"
const fileableType = "App\\Models\\JobVisit"

// CreateJobVisitRequest representa la solicitud para crear una visita
type CreateJobVisitRequest struct {
	UserID     int64   `json:"userId" example:"1"`
	ViewableBy *string `json:"viewableBy,omitempty" example:"[\"1\",\"2\"]"`
	Date       *string `json:"date,omitempty" example:"01-15-2024"`
	Report     *string `json:"report,omitempty" example:"Technician visited the site and inspected the HVAC system."`
}

// UpdateJobVisitRequest representa la solicitud para actualizar una visita
type UpdateJobVisitRequest struct {
	UserID     int64   `json:"userId" example:"1"`
	ViewableBy *string `json:"viewableBy,omitempty" example:"[\"1\",\"2\"]"`
	Date       *string `json:"date,omitempty" example:"01-15-2024"`
	Report     *string `json:"report,omitempty" example:"Updated report after follow-up visit."`
}

// JobVisitResponse representa la respuesta de una visita
type JobVisitResponse struct {
	ID         int64          `json:"id" example:"1"`
	JobID      int64          `json:"jobId" example:"1"`
	UserID     int64          `json:"userId" example:"1"`
	ViewableBy *string        `json:"viewableBy,omitempty"`
	Date       string         `json:"date" example:"01-15-2024"`
	Report     *string        `json:"report,omitempty"`
	Files      []FileResponse `json:"files,omitempty"`
	CreatedAt  string         `json:"createdAt,omitempty"`
	UpdatedAt  string         `json:"updatedAt,omitempty"`
}

// FileResponse representa la respuesta de un archivo
type FileResponse struct {
	ID        int64  `json:"id" example:"1"`
	Type      string `json:"type" example:"image"`
	URL       string `json:"url" example:"https://bucket.s3.amazonaws.com/uploads/1234_photo.jpg"`
	CreatedAt string `json:"createdAt,omitempty"`
}

func toVisitResponse(jv *domainJV.JobVisit) JobVisitResponse {
	resp := JobVisitResponse{
		ID:         jv.ID,
		JobID:      jv.JobID,
		UserID:     jv.UserID,
		ViewableBy: jv.ViewableBy,
		Date:       jv.Date.Format(dateFormat),
		Report:     jv.Report,
	}
	if jv.CreatedAt != nil {
		resp.CreatedAt = jv.CreatedAt.Format(timeFormat)
	}
	if jv.UpdatedAt != nil {
		resp.UpdatedAt = jv.UpdatedAt.Format(timeFormat)
	}
	return resp
}

func toFileResponse(f *domainFile.File) FileResponse {
	resp := FileResponse{
		ID:   f.ID,
		Type: f.GetDisplayType(),
		URL:  f.URL,
	}
	if f.CreatedAt != nil {
		resp.CreatedAt = f.CreatedAt.Format(timeFormat)
	}
	return resp
}

// List godoc
// @Summary      Listar visitas de un trabajo
// @Description  Obtiene una lista paginada de visitas de un trabajo específico
// @Tags         JobVisits
// @Accept       json
// @Produce      json
// @Param        jobId    path   int    true  "ID del trabajo"
// @Param        page     query  int    false "Número de página" default(1)
// @Param        pageSize query  int    false "Tamaño de página" default(10)
// @Param        search   query  string false "Búsqueda en reporte"
// @Param        userId   query  int    false "Filtrar por usuario"
// @Param        sort     query  string false "Campo de ordenamiento (date, created_at)" default(created_at)
// @Param        direction query string false "Dirección (ASC, DESC)" default(DESC)
// @Success      200 {object} response.PaginatedResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /jobs/{jobId}/visits [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.ParseInt(chi.URLParam(r, "jobId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de trabajo inválido")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}

	filters := map[string]interface{}{
		"jobId": jobID,
	}
	if search := r.URL.Query().Get("search"); search != "" {
		filters["search"] = search
	}
	if userID := r.URL.Query().Get("userId"); userID != "" {
		if uid, err := strconv.ParseInt(userID, 10, 64); err == nil {
			filters["userId"] = uid
		}
	}

	sort := r.URL.Query().Get("sort")
	direction := r.URL.Query().Get("direction")

	visits, total, err := h.useCase.List(r.Context(), filters, page, pageSize, sort, direction)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al obtener visitas")
		return
	}

	// Enrich con archivos
	var items []JobVisitResponse
	for _, v := range visits {
		resp := toVisitResponse(v)
		files, _ := h.fileUseCase.ListByFileable(r.Context(), v.ID, fileableType)
		for _, f := range files {
			resp.Files = append(resp.Files, toFileResponse(f))
		}
		items = append(items, resp)
	}

	response.Paginated(w, items, page, pageSize, int(total))
}

// Get godoc
// @Summary      Obtener visita por ID
// @Description  Obtiene una visita de trabajo por su ID
// @Tags         JobVisits
// @Accept       json
// @Produce      json
// @Param        jobId   path int true "ID del trabajo"
// @Param        visitId path int true "ID de la visita"
// @Success      200 {object} JobVisitResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /jobs/{jobId}/visits/{visitId} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	visitID, err := strconv.ParseInt(chi.URLParam(r, "visitId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de visita inválido")
		return
	}

	visit, err := h.useCase.GetByID(r.Context(), visitID)
	if err != nil {
		if err == domainJV.ErrJobVisitNotFound {
			response.Error(w, http.StatusNotFound, "Visita no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener visita")
		return
	}

	resp := toVisitResponse(visit)
	files, _ := h.fileUseCase.ListByFileable(r.Context(), visit.ID, fileableType)
	for _, f := range files {
		resp.Files = append(resp.Files, toFileResponse(f))
	}

	response.JSON(w, http.StatusOK, resp)
}

// Create godoc
// @Summary      Crear visita de trabajo
// @Description  Crea una nueva visita para un trabajo específico
// @Tags         JobVisits
// @Accept       json
// @Produce      json
// @Param        jobId path int true "ID del trabajo"
// @Param        request body CreateJobVisitRequest true "Datos de la visita"
// @Success      201 {object} JobVisitResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /jobs/{jobId}/visits [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.ParseInt(chi.URLParam(r, "jobId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de trabajo inválido")
		return
	}

	var req CreateJobVisitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	visit := &domainJV.JobVisit{
		JobID:      jobID,
		UserID:     req.UserID,
		ViewableBy: req.ViewableBy,
		Report:     req.Report,
		Date:       time.Now(),
	}

	if req.Date != nil {
		parsed, err := time.Parse(dateFormat, *req.Date)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Formato de fecha inválido, use MM-DD-YYYY")
			return
		}
		visit.Date = parsed
	}

	id, err := h.useCase.Create(r.Context(), visit)
	if err != nil {
		switch err {
		case domainJV.ErrJobNotFound:
			response.Error(w, http.StatusBadRequest, "Trabajo no encontrado")
		case domainJV.ErrUserNotFound:
			response.Error(w, http.StatusBadRequest, "Usuario no encontrado")
		default:
			response.Error(w, http.StatusInternalServerError, "Error al crear visita")
		}
		return
	}

	created, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al obtener visita creada")
		return
	}

	response.JSON(w, http.StatusCreated, toVisitResponse(created))
}

// Update godoc
// @Summary      Actualizar visita de trabajo
// @Description  Actualiza una visita de trabajo existente
// @Tags         JobVisits
// @Accept       json
// @Produce      json
// @Param        jobId   path int true "ID del trabajo"
// @Param        visitId path int true "ID de la visita"
// @Param        request body UpdateJobVisitRequest true "Datos actualizados"
// @Success      200 {object} JobVisitResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /jobs/{jobId}/visits/{visitId} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	visitID, err := strconv.ParseInt(chi.URLParam(r, "visitId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de visita inválido")
		return
	}

	var req UpdateJobVisitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	visit := &domainJV.JobVisit{
		ID:         visitID,
		UserID:     req.UserID,
		ViewableBy: req.ViewableBy,
		Report:     req.Report,
		Date:       time.Now(),
	}

	if req.Date != nil {
		parsed, err := time.Parse(dateFormat, *req.Date)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Formato de fecha inválido, use MM-DD-YYYY")
			return
		}
		visit.Date = parsed
	}

	if err := h.useCase.Update(r.Context(), visit); err != nil {
		switch err {
		case domainJV.ErrJobVisitNotFound:
			response.Error(w, http.StatusNotFound, "Visita no encontrada")
		case domainJV.ErrUserNotFound:
			response.Error(w, http.StatusBadRequest, "Usuario no encontrado")
		default:
			response.Error(w, http.StatusInternalServerError, "Error al actualizar visita")
		}
		return
	}

	updated, err := h.useCase.GetByID(r.Context(), visitID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al obtener visita actualizada")
		return
	}

	resp := toVisitResponse(updated)
	files, _ := h.fileUseCase.ListByFileable(r.Context(), updated.ID, fileableType)
	for _, f := range files {
		resp.Files = append(resp.Files, toFileResponse(f))
	}

	response.JSON(w, http.StatusOK, resp)
}

// Delete godoc
// @Summary      Eliminar visita de trabajo
// @Description  Elimina una visita de trabajo (soft delete). También elimina los archivos asociados.
// @Tags         JobVisits
// @Accept       json
// @Produce      json
// @Param        jobId   path int true "ID del trabajo"
// @Param        visitId path int true "ID de la visita"
// @Success      204
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /jobs/{jobId}/visits/{visitId} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	visitID, err := strconv.ParseInt(chi.URLParam(r, "visitId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de visita inválido")
		return
	}

	// Eliminar archivos asociados primero
	files, _ := h.fileUseCase.ListByFileable(r.Context(), visitID, fileableType)
	for _, f := range files {
		_ = h.fileUseCase.Delete(r.Context(), f.ID)
	}

	if err := h.useCase.Delete(r.Context(), visitID); err != nil {
		if err == domainJV.ErrJobVisitNotFound {
			response.Error(w, http.StatusNotFound, "Visita no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar visita")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListFiles godoc
// @Summary      Listar archivos de una visita
// @Description  Obtiene los archivos asociados a una visita de trabajo
// @Tags         JobVisits
// @Accept       json
// @Produce      json
// @Param        jobId   path int true "ID del trabajo"
// @Param        visitId path int true "ID de la visita"
// @Success      200 {array} FileResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /jobs/{jobId}/visits/{visitId}/files [get]
func (h *Handler) ListFiles(w http.ResponseWriter, r *http.Request) {
	visitID, err := strconv.ParseInt(chi.URLParam(r, "visitId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de visita inválido")
		return
	}

	files, err := h.fileUseCase.ListByFileable(r.Context(), visitID, fileableType)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al obtener archivos")
		return
	}

	var items []FileResponse
	for _, f := range files {
		items = append(items, toFileResponse(f))
	}

	response.JSON(w, http.StatusOK, items)
}

// UploadFile godoc
// @Summary      Subir archivo a una visita
// @Description  Sube uno o más archivos a una visita de trabajo (multipart/form-data, campo 'files')
// @Tags         JobVisits
// @Accept       multipart/form-data
// @Produce      json
// @Param        jobId   path     int    true "ID del trabajo"
// @Param        visitId path     int    true "ID de la visita"
// @Param        files   formData file   true "Archivos a subir (max 40MB c/u, max 20 archivos)"
// @Success      201 {array} FileResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /jobs/{jobId}/visits/{visitId}/files [post]
func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	visitID, err := strconv.ParseInt(chi.URLParam(r, "visitId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de visita inválido")
		return
	}

	// Verificar que la visita existe
	_, err = h.useCase.GetByID(r.Context(), visitID)
	if err != nil {
		if err == domainJV.ErrJobVisitNotFound {
			response.Error(w, http.StatusNotFound, "Visita no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al verificar visita")
		return
	}

	// Parsear multipart form (max 40MB por archivo * 20 archivos = ~800MB max)
	if err := r.ParseMultipartForm(800 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al parsear el formulario multipart")
		return
	}

	uploadedFiles := r.MultipartForm.File["files"]
	if len(uploadedFiles) == 0 {
		response.Error(w, http.StatusBadRequest, "No se enviaron archivos")
		return
	}
	if len(uploadedFiles) > 20 {
		response.Error(w, http.StatusBadRequest, "Máximo 20 archivos por solicitud")
		return
	}

	var results []FileResponse
	for _, fileHeader := range uploadedFiles {
		// Validar tamaño (40MB max)
		if fileHeader.Size > 40<<20 {
			response.Error(w, http.StatusBadRequest, "El archivo "+fileHeader.Filename+" excede el tamaño máximo de 40MB")
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Error al abrir archivo: "+fileHeader.Filename)
			return
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		uploaded, err := h.fileUseCase.Upload(r.Context(), visitID, fileableType, fileHeader.Filename, contentType, file)
		_ = file.Close()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Error al subir archivo: "+fileHeader.Filename)
			return
		}

		results = append(results, toFileResponse(uploaded))
	}

	response.JSON(w, http.StatusCreated, results)
}

// DeleteFile godoc
// @Summary      Eliminar archivo de una visita
// @Description  Elimina un archivo de S3 y de la base de datos
// @Tags         JobVisits
// @Accept       json
// @Produce      json
// @Param        jobId   path int true "ID del trabajo"
// @Param        visitId path int true "ID de la visita"
// @Param        fileId  path int true "ID del archivo"
// @Success      204
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /jobs/{jobId}/visits/{visitId}/files/{fileId} [delete]
func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, "fileId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de archivo inválido")
		return
	}

	if err := h.fileUseCase.Delete(r.Context(), fileID); err != nil {
		if err == domainFile.ErrFileNotFound {
			response.Error(w, http.StatusNotFound, "Archivo no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar archivo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DownloadFile godoc
// @Summary      Descargar archivo de una visita
// @Description  Descarga un archivo de S3
// @Tags         JobVisits
// @Produce      octet-stream
// @Param        jobId   path int true "ID del trabajo"
// @Param        visitId path int true "ID de la visita"
// @Param        fileId  path int true "ID del archivo"
// @Success      200
// @Failure      400 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Router       /jobs/{jobId}/visits/{visitId}/files/{fileId}/download [get]
func (h *Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(chi.URLParam(r, "fileId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de archivo inválido")
		return
	}

	body, contentType, filename, err := h.fileUseCase.Download(r.Context(), fileID)
	if err != nil {
		if err == domainFile.ErrFileNotFound {
			response.Error(w, http.StatusNotFound, "Archivo no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al descargar archivo")
		return
	}
	defer func() { _ = body.Close() }()

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.WriteHeader(http.StatusOK)

	buf := make([]byte, 32*1024)
	for {
		n, readErr := body.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return
			}
		}
		if readErr != nil {
			break
		}
	}
}
