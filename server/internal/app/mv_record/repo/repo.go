package repo

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

func SetFieldValue(db *gorm.DB, recordId string, fieldKey string, value any) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("invalid value: %w", err)
	}

	sql := `
		UPDATE mv_record
		SET row_data = jsonb_set(row_data, ?, ?::jsonb, true),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	return db.Exec(sql, fmt.Sprintf("{%s}", fieldKey), string(valueJSON), recordId).Error
}
