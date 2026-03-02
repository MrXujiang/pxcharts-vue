package service

import (
	"encoding/json"

	"gorm.io/datatypes"
)

// MapToJSON 将 map[string]any 转换为 datatypes.JSON
func MapToJSON(m map[string]any) datatypes.JSON {
	if m == nil {
		return nil
	}

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return datatypes.JSON("{}")
	}

	return datatypes.JSON(jsonBytes)
}
