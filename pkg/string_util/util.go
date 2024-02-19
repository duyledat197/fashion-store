package stringutil

// Coalesce ...
func Coalesce(src ...string) string {
	for _, s := range src {
		if s != "" {
			return s
		}
	}

	return ""
}
