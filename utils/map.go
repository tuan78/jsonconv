package utils

func CopyMap(from map[string]interface{}, to map[string]interface{}) {
	for k, v := range from {
		to[k] = v
	}
}
