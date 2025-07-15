package btree

// Item ...
type Item[K, V any] struct {
	Key   K
	Value V
}

// NewItem ...
func NewItem[K, V any](k K, value V) *Item[K, V] {
	return &Item[K, V]{
		Key:   k,
		Value: value,
	}
}
