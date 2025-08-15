package testutils

// ApplyHooks применяет хуки.
func ApplyHooks[T any](value T, hooks []func(T)) {
	for _, hook := range hooks {
		hook(value)
	}
}
