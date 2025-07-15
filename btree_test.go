package btree

import (
	"cmp"
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

func TestBTreeEndToEnd(t *testing.T) {
	n := 1000

	values := rand.Perm(n)
	t.Log(values)
	tree := arrange(4, values)

	slices.Sort(values)
	count := tree.Count()

	for _, v := range rand.Perm(n) {
		t.Log(v)

		c := slices.Contains(values, v)

		if c {
			count--
			ind := slices.Index(values, v)
			values = slices.Delete(values, ind, ind+1)
		}
		r := tree.TryRemove(v)

		assert(t, tree, c, r, count, values)
	}
}

func TestTryAddRoot(t *testing.T) {
	tree := arrange(2, []int{})
	value := 1
	expected := true
	expCount := 1
	expContent := []int{value}

	actual := tree.TryAdd(value, value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryAddExisting(t *testing.T) {
	tree := arrange(2, []int{1})
	value := 1
	expected := false
	expCount := 1
	expContent := []int{value}

	actual := tree.TryAdd(value, value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryAddInVacantNode(t *testing.T) {
	tree := arrange(2, []int{1})
	value := 2
	expected := true
	expCount := 2
	expContent := []int{1, value}

	actual := tree.TryAdd(value, value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryAddSplitRoot(t *testing.T) {
	tree := arrange(2, []int{7, 10, 15})
	value := 8
	expected := true
	expCount := 4
	expContent := []int{7, value, 10, 15}

	actual := tree.TryAdd(value, value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryAddSplitWithChildren(t *testing.T) {
	tree := arrange(2, []int{10, 15, 7, 8, 20, 3, 18, 0, 12})
	value := 9
	expected := true
	expCount := 10
	expContent := []int{0, 3, 7, 8, value, 10, 12, 15, 18, 20}

	actual := tree.TryAdd(value, value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryAddSplitLeaf(t *testing.T) {
	tree := arrange(2, []int{10, 15, 7, 8, 20, 3})
	value := 0
	expected := true
	expCount := 7
	expContent := []int{value, 3, 7, 8, 10, 15, 20}

	actual := tree.TryAdd(value, value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryAddLeftSplit(t *testing.T) {
	tree := arrange(2, []int{7, 10, 15})
	value := 0
	expected := true
	expCount := 4
	expContent := []int{value, 7, 10, 15}

	actual := tree.TryAdd(value, value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryAddRightSplit(t *testing.T) {
	tree := arrange(2, []int{7, 10, 15})
	value := 12
	expected := true
	expCount := 4
	expContent := []int{7, 10, value, 15}

	actual := tree.TryAdd(value, value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryGetNilRoot(t *testing.T) {
	tree := arrange(2, []int{})
	value := 0
	expected := false

	actual, v := tree.TryGet(value)

	if expected != actual && v == value {
		t.Errorf(`
		Expected: %t value %d
		Actual:   %t value %d`, expected, value, actual, v)
	}
}

func TestTryGetNonExisting(t *testing.T) {
	value := 0
	tree := arrange(2, []int{1, 2, 3, 4, 5, 6, 7})
	expected := false

	actual, v := tree.TryGet(value)

	if expected != actual && v == value {
		t.Errorf(`
		Expected: %t value %d
		Actual:   %t value %d`, expected, value, actual, v)
	}
}

func TestTryGetExisting(t *testing.T) {
	value := 6
	tree := arrange(2, []int{1, 2, 3, 4, 5, value})
	expected := true

	actual, v := tree.TryGet(value)

	if expected != actual && v == value {
		t.Errorf(`
		Expected: %t value %d
		Actual:   %t value %d`, expected, value, actual, v)
	}
}

func TestTryRemoveNonExistLeafRoot(t *testing.T) {
	tree := arrange(2, []int{1, 2, 3})
	value := 4
	expected := false
	expCount := 3
	expContent := []int{1, 2, 3}

	actual := tree.TryRemove(value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryRemoveNonExistNonLeafRoot(t *testing.T) {
	tree := arrange(2, []int{0, 10, 20, 30, 40, 50, 60, 15, 25, 27, 70, 80, 90})
	value := 29
	expected := false
	expCount := 13
	expContent := []int{0, 10, 15, 20, 25, 27, 30, 40, 50, 60, 70, 80, 90}

	actual := tree.TryRemove(value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryRemoveNullRoot(t *testing.T) {
	tree := arrange(2, []int{})
	value := 1
	expected := false
	expCount := 0
	expContent := []int{}

	actual := tree.TryRemove(value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryRemoveLeafRoot(t *testing.T) {
	tree := arrange(2, []int{1, 2, 3})
	value := 2
	expected := true
	expCount := 2
	expContent := []int{1, 3}

	actual := tree.TryRemove(value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryRemoveLastRootKey(t *testing.T) {
	tree := arrange(2, []int{1})
	value := 1
	expected := true
	expCount := 0
	expContent := []int{}

	actual := tree.TryRemove(value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryRemoveMergeRoot(t *testing.T) {
	tree := arrange(2, []int{1, 2, 3, 4})
	value := 2
	expected := true
	expCount := 2
	expContent := []int{1, 3}

	tree.TryRemove(4)
	actual := tree.TryRemove(value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryRemoveLeafStealLeft(t *testing.T) {
	tree := arrange(2, []int{1, 2, 3, 0})
	value := 3
	expected := true
	expCount := 3
	expContent := []int{0, 1, 2}

	actual := tree.TryRemove(value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryRemoveLeafStealRight(t *testing.T) {
	tree := arrange(2, []int{1, 2, 3, 4})
	value := 1
	expected := true
	expCount := 3
	expContent := []int{2, 3, 4}

	actual := tree.TryRemove(value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryRemoveLeafLeftMerge(t *testing.T) {
	tree := arrange(2, []int{3, 4, 5, 1, 2, 0})
	value := 5
	expected := true
	expCount := 5
	expContent := []int{0, 1, 2, 3, 4}

	actual := tree.TryRemove(value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryRemoveLeafRightMerge(t *testing.T) {
	tree := arrange(2, []int{1, 2, 3, 4, 5, 6})
	value := 1
	expected := true
	expCount := 5
	expContent := []int{2, 3, 4, 5, 6}

	actual := tree.TryRemove(value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryRemoveInternalPredecessor(t *testing.T) {
	tree := arrange(2, []int{0, 10, 20, 30, 40, 50, 60, 15, 25, 27, 70, 80, 90})
	value := 30
	expected := true
	expCount := 12
	expContent := []int{0, 10, 15, 20, 25, 27, 40, 50, 60, 70, 80, 90}

	actual := tree.TryRemove(value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryRemoveInternalSuccessor(t *testing.T) {
	tree := arrange(2, []int{0, 10, 20, 30, 40, 50, 60, 15, 25, 27, 70, 80, 90})
	value := 20
	expected := true
	expCount := 12
	expContent := []int{0, 10, 15, 25, 27, 30, 40, 50, 60, 70, 80, 90}

	actual := tree.TryRemove(value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func TestTryRemoveInternalMerge(t *testing.T) {
	tree := arrange(2, []int{0, 10, 20, 30, 40, 50, 60, 15, 25, 27, 70, 80, 90})
	value := 10
	expected := true
	expCount := 12
	expContent := []int{0, 15, 20, 25, 27, 30, 40, 50, 60, 70, 80, 90}

	actual := tree.TryRemove(value)

	assert(t, tree, expected, actual, uint64(expCount), expContent)
}

func arrange[K cmp.Ordered](t int, s []K) *BTree[K, K] {
	tree := New[K, K](t)
	for _, v := range s {
		tree.TryAdd(v, v)
	}
	return tree
}

func assert(t *testing.T, tree *BTree[int, int], expected, actual bool, expCount uint64, expContent []int) {
	inorder := tree.InorderTraverse()
	inv, mes := tree.checkInvariant()
	if actual != expected ||
		tree.Count() != expCount ||
		!inv ||
		!equal(inorder, expContent) {
		t.Errorf(`
		Expected: result %t count %d content %v
		Actual:   result %t count %d content %v invariant = %s`,
			expected, expCount, expContent,
			actual, tree.Count(), inorder, mes)
	}
}

func equal[K cmp.Ordered, V any](items []*Item[K, V], keys []K) bool {
	if len(items) != len(keys) {
		return false
	}
	for i, v := range items {
		if v.Key != keys[i] {
			return false
		}
	}
	return true
}

// Every node contains at least n - 1 keys beside root node.
// Every node contains less or equal than 2t - 1 keys and 2t children.
// Keys in accend order.
// Every leaf on same height.
func (b BTree[K, V]) checkInvariant() (bool, string) {
	if b.root == nil {
		return true, ""
	}
	curr := b.root
	if len(curr.items) > 2*b.t-1 {
		return false, fmt.Sprintf("Root: %v, len(curr.keys) > 2*t-1", curr.items)
	}
	if len(curr.children) > 2*b.t {
		return false, "Root children len > 2 * t"
	}

	if len(curr.children) != 0 {
		var ok, height, mes = b._checkInvariant(curr.children[0], 0)
		if !ok {
			return false, mes
		}
		for _, v := range curr.children[1:] {
			ok, childH, mes := b._checkInvariant(v, 0)

			if !ok {
				return false, mes
			}

			if height != childH {
				return false, fmt.Sprintf("Keys: %v, children heights are different", curr.items)
			}
		}
	}
	return true, ""
}

func (b BTree[K, V]) _checkInvariant(curr *node[K, V], h int) (bool, int, string) {
	h++
	if len(curr.items) < b.t-1 || len(curr.items) > 2*b.t-1 {
		return false, h, fmt.Sprintf("Keys: %v, len(curr.keys) < t-1 or len(curr.keys) > 2*t-1", curr.items)
	}
	if len(curr.children) > 2*b.t {
		return false, h, fmt.Sprintf("Keys: %v, len(curr.children) > 2*t", curr.items)
	}

	if len(curr.children) != 0 {
		var ok, height, mes = b._checkInvariant(curr.children[0], h)

		if !ok {
			return false, h, mes
		}

		for _, v := range curr.children[1:] {
			ok, childH, mes := b._checkInvariant(v, h)
			if !ok {
				return false, h, mes
			}

			if height != childH {
				return false, h, fmt.Sprintf("Keys: %v, children heights are different", curr.items)
			}
		}
	}
	return true, h, ""
}
