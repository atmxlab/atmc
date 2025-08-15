package orderedset

import "fmt"

// OrderedSet упорядоченная map
type OrderedSet[K comparable, V any] struct {
	orderedKeys []K
	m           map[K]V
}

func New[K comparable, V any](capacity int) *OrderedSet[K, V] {
	return &OrderedSet[K, V]{
		orderedKeys: make([]K, 0, capacity),
		m:           make(map[K]V, capacity),
	}
}

// Values получение всех значений.
func (os *OrderedSet[K, V]) Values() []V {
	values := make([]V, 0, len(os.orderedKeys))
	for _, v := range os.Iterator() {
		values = append(values, v)
	}

	return values
}

func (os *OrderedSet[K, V]) Keys() []K {
	keys := make([]K, 0, len(os.orderedKeys))
	for k, _ := range os.Iterator() {
		keys = append(keys, k)
	}

	return keys
}

// Len получение длины.
func (os *OrderedSet[K, V]) Len() int {
	return len(os.orderedKeys)
}

// Iterator генератор, который позволяет итерироваться по map в правильном порядке и в конструкции for range.
func (os *OrderedSet[K, V]) Iterator() func(yield func(k K, v V) bool) {
	return func(yield func(K, V) bool) {
		for _, k := range os.orderedKeys {
			v := os.m[k]
			if !yield(k, v) {
				return
			}
		}
	}
}

func (os *OrderedSet[K, V]) Set(k K, v V) {
	if _, ok := os.m[k]; !ok {
		os.orderedKeys = append(os.orderedKeys, k)
	}

	os.m[k] = v
}

func (os *OrderedSet[K, V]) Get(k K) (V, bool) {
	v, ok := os.m[k]
	return v, ok
}

// GetValue значение не обязано быть, если нет, будет zero value.
func (os *OrderedSet[K, V]) GetValue(k K) V {
	v, _ := os.Get(k)
	return v
}

// MustGet значение обязано быть, если нет будет паника.
func (os *OrderedSet[K, V]) MustGet(k K) V {
	if v, ok := os.Get(k); ok {
		return v
	}

	panic(fmt.Sprintf("missing key: %v", k))
}

func (os *OrderedSet[K, V]) Delete(k K) {
	if _, ok := os.m[k]; !ok {
		return
	}

	os.deleteFromOrderedKey(k)
	os.deleteFromMap(k)
}

func (os *OrderedSet[K, V]) deleteFromMap(k K) {
	delete(os.m, k)
}

func (os *OrderedSet[K, V]) deleteFromOrderedKey(k K) {
	idx := os.findKeyIndex(k)

	if idx < 0 || idx >= len(os.orderedKeys) {
		return
	}

	if idx == len(os.orderedKeys)-1 {
		os.orderedKeys = os.orderedKeys[:len(os.orderedKeys)-1]
	} else {
		os.orderedKeys = append(os.orderedKeys[:idx], os.orderedKeys[idx+1:]...)
	}
}

func (os *OrderedSet[K, V]) findKeyIndex(k K) int {
	for i, key := range os.orderedKeys {
		if k == key {
			return i
		}
	}

	return -1
}
