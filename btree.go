package btree

import (
	"cmp"
	"slices"
)

// BTree ...
type BTree[K, V any] struct {
	root  *node[K, V]
	t     int
	count uint64

	compare func(K, K) int
}

// New ...
func New[K cmp.Ordered, V any](t int) *BTree[K, V] {
	return &BTree[K, V]{
		t:       t,
		compare: cmp.Compare[K],
	}
}

// NewWith ...
func NewWith[K, V any](t int, compare func(K, K) int) *BTree[K, V] {
	return &BTree[K, V]{
		t:       t,
		compare: compare,
	}
}

func (b *BTree[K, V]) Count() uint64 {
	return b.count
}

// TryAdd ...
func (b *BTree[K, V]) TryAdd(key K, value V) bool {
	if b.root == nil {
		b.root = newNode[K, V](b.t)
		b.root.items = append(b.root.items, NewItem(key, value))
		b.count++
		return true
	}

	curr := b.root
	var parent *node[K, V] = nil

	for {
		ind, exist := curr.bSearch(key, b.compare)
		if exist {
			return false
		}
		if curr.isFilled() {
			left, right := b.splitNode(curr, parent)
			if ind <= b.t-1 {
				curr = left
			} else {
				curr = right
			}
			continue
		}
		if curr.isLeaf() {
			curr.items = slices.Insert(curr.items, ind, NewItem(key, value))
			b.count++
			return true
		}

		parent = curr
		curr = curr.children[ind]
	}
}

// TryGet ...
func (b *BTree[K, V]) TryGet(key K) (ok bool, value V) {
	node := b.root
	for node != nil {
		ind, ok := node.bSearch(key, b.compare)
		if ok {
			return true, node.items[ind].Value
		}
		if node.isLeaf() {
			return false, value
		}
		node = node.children[ind]
	}
	return false, value
}

func (b *BTree[K, V]) TryUpdate(key K, new V) (ok bool) {
	node := b.root
	for node != nil {
		ind, ok := node.bSearch(key, b.compare)
		if ok {
			node.items[ind].Value = new
			return true
		}
		if node.isLeaf() {
			return false
		}
		node = node.children[ind]
	}
	return false
}

// TryRemove ...
func (b *BTree[K, V]) TryRemove(key K) bool {
	if b.root == nil {
		return false
	}

	ind, exist := b.root.bSearch(key, b.compare)

	if exist {
		b.delete(b.root, ind)

		if len(b.root.items) == 0 {
			if b.root.isLeaf() {
				b.root = nil
			} else {
				b.root = b.root.children[0]
			}
		}
		b.count--
		return true
	}
	if b.root.isLeaf() {
		return false
	}

	removed := b.remove(b.root.children[ind], b.root, ind, key)
	if len(b.root.items) == 0 {
		b.root = b.root.children[0]
	}
	if removed {
		b.count--
		return true
	}
	return false
}

func (b *BTree[K, V]) InorderTraverse() []Item[K, V] {
	if b.root == nil {
		return []Item[K, V]{}
	}
	return inorderTraverse(b.root)
}

func (b *BTree[K, V]) remove(node, parent *node[K, V], childInd int, key K) bool {
	for {
		if len(node.items) == b.t-1 {
			if !b.stealIfPossible(parent, childInd) {
				node = b.merge(parent, childInd)
			}
		}

		ind, exist := node.bSearch(key, b.compare)
		if exist {
			b.delete(node, ind)
			return true
		}

		if node.isLeaf() {
			return false
		}

		node, parent, childInd = node.children[ind], node, ind
	}
}

func (b *BTree[K, V]) delete(node *node[K, V], valueInd int) {
	if node.isLeaf() {
		node.items = slices.Delete(node.items, valueInd, valueInd+1)
	} else {
		if len(node.children[valueInd].items) > b.t-1 {
			m := b.max(node.children[valueInd])

			b.remove(node.children[valueInd], node, valueInd, m.Key)

			node.items[valueInd] = m
		} else if len(node.children[valueInd+1].items) > b.t-1 {
			m := b.min(node.children[valueInd+1])

			b.remove(node.children[valueInd+1], node, valueInd+1, m.Key)

			node.items[valueInd] = m
		} else {
			leftLen := len(node.children[valueInd].items)
			b.merge(node, valueInd+1)
			b.delete(node.children[valueInd], leftLen)
		}
	}
}

func (b *BTree[K, V]) stealIfPossible(parent *node[K, V], childInd int) bool {
	child := parent.children[childInd]
	if childInd > 0 && len(parent.children[childInd-1].items) >= b.t { // steal from left sibling
		sibling := parent.children[childInd-1]

		child.items = slices.Insert(child.items, 0, parent.items[childInd-1])

		if !sibling.isLeaf() {
			child.children = slices.Insert(child.children, 0, sibling.children[len(sibling.children)-1])
			sibling.children = slices.Delete(sibling.children, len(sibling.children)-1, len(sibling.children))
		}

		parent.items[childInd-1] = sibling.items[len(sibling.items)-1]

		sibling.items = slices.Delete(sibling.items, len(sibling.items)-1, len(sibling.items))

	} else if childInd+1 < len(parent.children) && len(parent.children[childInd+1].items) >= b.t {
		sibling := parent.children[childInd+1]

		child.items = append(child.items, parent.items[childInd])

		if !sibling.isLeaf() {
			child.children = append(child.children, sibling.children[0])
			sibling.children = slices.Delete(sibling.children, 0, 1)
		}

		parent.items[childInd] = sibling.items[0]
		sibling.items = slices.Delete(sibling.items, 0, 1)
	} else {
		return false
	}
	return true
}

func (b *BTree[K, V]) splitNode(node, parent *node[K, V]) (left, right *node[K, V]) {
	left = newNode[K, V](b.t)
	left.items = node.items[:b.t-1]

	item := node.items[b.t-1]

	right = newNode[K, V](b.t)
	right.items = append(right.items, node.items[b.t:]...)

	if !node.isLeaf() {
		left.children = node.children[:b.t]
		right.children = append(right.children, node.children[b.t:]...)
	}

	if parent == nil {
		parent = newNode[K, V](b.t)
		parent.items = append(parent.items, item)
		parent.children = append(parent.children, left, right)

		b.root = parent
	} else {
		ind, _ := parent.bSearch(item.Key, b.compare)

		parent.items = slices.Insert(parent.items, ind, item)
		parent.children[ind] = left
		parent.children = slices.Insert(parent.children, ind+1, right)
	}
	return
}

func (b *BTree[K, V]) merge(parent *node[K, V], childInd int) *node[K, V] {
	// parentInd = left
	var left, right int = childInd, childInd
	if childInd > 0 {
		left--
	} else {
		right++
	}

	leftSibl := parent.children[left]
	leftSibl.items = append(append(leftSibl.items,
		parent.items[left]),
		parent.children[right].items...)
	leftSibl.children = append(leftSibl.children, parent.children[right].children...)

	parent.items = slices.Delete(parent.items, left, right)
	parent.children = slices.Delete(parent.children, right, right+1)

	return leftSibl
}

func inorderTraverse[K, V any](node *node[K, V]) []Item[K, V] {
	result := []Item[K, V]{}
	if node.isLeaf() {
		for _, v := range node.items {
			result = append(result, *v)
		}
		return result
	}
	for i, v := range node.items {
		result = append(result, inorderTraverse(node.children[i])...)
		result = append(result, *v)
	}
	result = append(result, inorderTraverse(node.children[len(node.children)-1])...)

	return result
}

func (b BTree[K, V]) min(node *node[K, V]) *Item[K, V] {
	for !node.isLeaf() {
		node = node.children[0]
	}
	return node.items[0]
}

func (b BTree[K, V]) max(node *node[K, V]) *Item[K, V] {
	for !node.isLeaf() {
		node = node.children[len(node.children)-1]
	}
	return node.items[len(node.items)-1]
}
