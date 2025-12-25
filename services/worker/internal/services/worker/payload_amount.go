package worker

import "encoding/json"

func extractAmount(payload json.RawMessage) float64 {
	if len(payload) == 0 {
		return 0
	}

	var m map[string]any
	if err := json.Unmarshal(payload, &m); err != nil {
		return 0
	}

	if v, ok := m["amount"]; ok {
		switch t := v.(type) {
		case float64:
			return t
		case int:
			return float64(t)
		case json.Number:
			f, _ := t.Float64()
			return f
		case string:
			var num json.Number = json.Number(t)
			f, _ := num.Float64()
			return f
		}
	}

	items, ok := m["items"].([]any)
	if !ok {
		return 0
	}

	var total float64
	for _, it := range items {
		im, ok := it.(map[string]any)
		if !ok {
			continue
		}

		price := toFloat(im["price"])
		qty := toFloat(im["qty"])
		if qty == 0 {
			qty = 1
		}
		total += price * qty
	}
	return total
}

func toFloat(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case json.Number:
		f, _ := t.Float64()
		return f
	case string:
		var num json.Number = json.Number(t)
		f, _ := num.Float64()
		return f
	default:
		return 0
	}
}
