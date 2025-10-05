package orderedmap

type OrderedMap[K comparable, V any] struct {
	keys   []K
	values map[K]V
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		keys:   make([]K, 0),
		values: make(map[K]V),
	}
}

func (om *OrderedMap[K, V]) Set(key K, value V) {
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	val, exists := om.values[key]
	return val, exists
}

func (om *OrderedMap[K, V]) Delete(key K) bool {
	if _, exists := om.values[key]; exists {
		delete(om.values, key)
		for i, k := range om.keys {
			if k == key {
				om.keys = append(om.keys[:i], om.keys[i+1:]...)
				break
			}
		}
		return true
	}
	return false
}

func (om *OrderedMap[K, V]) Keys() []K {
	return om.keys
}

func (om *OrderedMap[K, V]) Values() []V {
	values := make([]V, len(om.keys))
	for i, key := range om.keys {
		values[i] = om.values[key]
	}
	return values
}

func (om *OrderedMap[K, V]) Entries() ([]K, []V) {
	values := make([]V, len(om.keys))
	for i, key := range om.keys {
		values[i] = om.values[key]
	}
	return om.keys, values
}
