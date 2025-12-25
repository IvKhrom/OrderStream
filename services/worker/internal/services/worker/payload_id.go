package worker

import (
	"encoding/json"
	"fmt"
)

func payloadID(payload json.RawMessage) string {
	var tmp map[string]any
	if err := json.Unmarshal(payload, &tmp); err != nil {
		return ""
	}
	if v, ok := tmp["id"]; ok {
		switch t := v.(type) {
		case string:
			return t
		case float64:
			return fmt.Sprintf("%v", t)
		default:
			return ""
		}
	}
	return ""
}
