package pointer

func To[T any](v T) *T {
	return &v
}

func From[T any](p *T) T {
	if p == nil {
		var defaultValue T
		return defaultValue
	}
	return *p
}
