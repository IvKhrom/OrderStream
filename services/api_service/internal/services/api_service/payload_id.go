package api_service

import "encoding/json"

func payloadID(payload json.RawMessage) string {
	var tmp map[string]any
	if err := json.Unmarshal(payload, &tmp); err != nil {
		return ""
	}
	if v, ok := tmp["id"].(string); ok {
		return v
	}
	return ""
}


