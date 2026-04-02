package job

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	domainJobSMS "github.com/your-org/jvairv2/pkg/domain/job_sms"
	handler "github.com/your-org/jvairv2/pkg/rest/handler"
)

// SendDispatchSMSRequest representa la solicitud para enviar SMS de dispatch
type SendDispatchSMSRequest struct {
	PhoneNumbers string `json:"phoneNumbers"`
	TextMessage  string `json:"textMessage"`
}

// SendDispatchSMS envía SMS de dispatch del job a uno o más números de teléfono
// @Summary Enviar SMS de dispatch
// @Description Envía SMS de dispatch a uno o más números de teléfono via AWS SNS o Twilio
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path int true "Job ID"
// @Param request body SendDispatchSMSRequest true "SMS request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/{id}/dispatch-sms [post]
func (h *Handler) SendDispatchSMS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jobIDStr := chi.URLParam(r, "id")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "ID de job inválido")
		return
	}

	var req SendDispatchSMSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Request inválido")
		return
	}

	if strings.TrimSpace(req.PhoneNumbers) == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "phoneNumbers es requerido")
		return
	}
	if strings.TrimSpace(req.TextMessage) == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "textMessage es requerido")
		return
	}

	if h.smsSender == nil {
		handler.RespondWithError(w, http.StatusInternalServerError, "Servicio de SMS no configurado")
		return
	}

	nonDigits := regexp.MustCompile(`[^0-9]`)
	rawNumbers := strings.Split(req.PhoneNumbers, ",")

	var failedNumbers []string
	var sentNumbers []string

	slog.InfoContext(ctx, "[SMS dispatch] iniciando envío",
		slog.Int64("jobId", jobID),
		slog.String("phoneNumbers", req.PhoneNumbers),
	)

	for _, raw := range rawNumbers {
		phone := strings.TrimSpace(raw)
		if phone == "" {
			continue
		}

		justNums := nonDigits.ReplaceAllString(phone, "")

		slog.DebugContext(ctx, "[SMS dispatch] normalizando número",
			slog.String("original", phone),
			slog.String("soloDigitos", justNums),
			slog.Int("longitud", len(justNums)),
		)

		// Eliminar 1 inicial si tiene 11 dígitos
		if len(justNums) == 11 && strings.HasPrefix(justNums, "1") {
			justNums = justNums[1:]
			slog.DebugContext(ctx, "[SMS dispatch] eliminado prefijo 1", slog.String("resultado", justNums))
		}

		// Validar que queden exactamente 10 dígitos
		if len(justNums) != 10 {
			slog.WarnContext(ctx, "[SMS dispatch] número inválido (no tiene 10 dígitos)",
				slog.String("original", phone),
				slog.String("soloDigitos", justNums),
				slog.Int("longitud", len(justNums)),
			)
			failedNumbers = append(failedNumbers, phone)
			continue
		}

		e164 := "+1" + justNums
		slog.InfoContext(ctx, "[SMS dispatch] enviando SMS",
			slog.String("e164", e164),
		)

		if err := h.smsSender.SendSMS(ctx, e164, req.TextMessage); err != nil {
			slog.ErrorContext(ctx, "[SMS dispatch] fallo al enviar SMS",
				slog.String("e164", e164),
				slog.String("original", phone),
				slog.String("error", err.Error()),
			)
			failedNumbers = append(failedNumbers, phone)
			continue
		}

		slog.InfoContext(ctx, "[SMS dispatch] SMS enviado correctamente", slog.String("e164", e164))
		sentNumbers = append(sentNumbers, phone)
	}

	if len(failedNumbers) > 0 {
		handler.RespondWithError(w, http.StatusBadRequest,
			fmt.Sprintf("SMS no pudo enviarse a los siguientes números: %s", strings.Join(failedNumbers, ", ")))
		return
	}

	// Registrar el SMS enviado
	if h.jobSMSService != nil && len(sentNumbers) > 0 {
		smsRecord := &domainJobSMS.JobSMS{
			JobID:      jobID,
			Recipients: sentNumbers,
			Type:       "sms sent",
			Message:    req.TextMessage,
		}
		_ = h.jobSMSService.Create(ctx, smsRecord)
	}

	// Registrar actividad
	h.logActivity(ctx, jobID, "job_sms_dispatched", "Job was dispatched via sms")

	// Responder con éxito
	handler.RespondWithSuccess(w, fmt.Sprintf("SMS enviado a %s", req.PhoneNumbers))
}
