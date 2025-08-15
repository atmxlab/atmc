package processor

type importStack map[string]struct{}

func newEmptyImportStack() importStack {
	return make(importStack)
}

func (is importStack) Clone() importStack {
	copied := make(map[string]struct{}, len(is))
	for path, v := range is {
		copied[path] = v
	}

	return copied
}
