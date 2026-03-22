package job_sms

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/job_sms"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) job_sms.Repository {
	return &repository{db: db}
}

func encodeRecipients(recipients []string) ([]byte, error) {
	return json.Marshal(recipients)
}

func decodeRecipients(raw []byte) []string {
	if len(raw) == 0 {
		return []string{}
	}

	var recipients []string
	if err := json.Unmarshal(raw, &recipients); err == nil {
		return recipients
	}

	legacy := strings.Split(string(raw), ",")
	result := make([]string, 0, len(legacy))
	for _, item := range legacy {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
