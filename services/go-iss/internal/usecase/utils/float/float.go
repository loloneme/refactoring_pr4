package float

import "strconv"

func ExtractFloat(m map[string]interface{}, key string) *float64 {
	if x, ok := m[key]; ok {
		if f, ok := x.(float64); ok {
			return &f
		}
		if s, ok := x.(string); ok {
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				return &f
			}
		}
	}
	return nil
}

func ParseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
