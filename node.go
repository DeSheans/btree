package btree

type node[K, V any] struct {
	items    []*Item[K, V]
	children []*node[K, V]
}

func newNode[K, V any](t int) *node[K, V] {
	return &node[K, V]{
		items:    make([]*Item[K, V], 0, (t<<1)-1),
		children: make([]*node[K, V], 0, (t << 1)),
	}
}

func (n *node[K, V]) isLeaf() bool {
	return len(n.children) == 0
}

func (n *node[K, V]) isFilled() bool {
	return len(n.items) == cap(n.items)
}

func (n node[K, V]) bSearch(key K, cmp func(K, K) int) (int, bool) {
	l := len(n.items)
	i, j := 0, l
	for i < j {
		h := int(uint(i+j) >> 1)
		if cmp(n.items[h].Key, key) < 0 {
			i = h + 1
		} else {
			j = h
		}
	}
	return i, i < l && cmp(n.items[i].Key, key) == 0
}

