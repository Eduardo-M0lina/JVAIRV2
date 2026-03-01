package warranty_equipment

import (
	"database/sql"
	"time"
)

func assignNullString(target **string, ns sql.NullString) {
	if ns.Valid {
		*target = &ns.String
	}
}

func assignNullTime(target **time.Time, nt sql.NullTime) {
	if nt.Valid {
		*target = &nt.Time
	}
}
