package filter

func RemoveFieldFilter(fields []string) FilterFunc {
	return func(m map[string]any) map[string]any {
		for _, field := range fields {
			delete(m, field)
		}
		return m
	}
}
