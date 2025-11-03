package metamysql

import (
	"database/sql"
	"time"
)

func NullStringToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func NullTimeToTime(t sql.NullTime) time.Time {
	if t.Valid {
		return t.Time
	}
	return time.Time{}
}
