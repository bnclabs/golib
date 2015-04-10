// A Left-Leaning Red-Black (LLRB) implementation of 2-3 balanced
// binary search trees, based on the following work:
//
//   http://www.cs.princeton.edu/~rs/talks/LLRB/08Penn.pdf
//   http://www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf
//   http://www.cs.princeton.edu/~rs/talks/LLRB/Java/RedBlackBST.java
//
// 2-3 trees (and the run-time equivalent 2-3-4 trees) are the de
// facto standard BST //  algoritms found in implementations of Python,
// Java, and other libraries. The LLRB implementation of 2-3 trees
// is a recent improvement on the traditional implementation,
// observed and documented by Robert Sedgewick.
package llrb

// LLRB is a Left-Leaning Red-Black (LLRB) implementation
// of 2-3 trees
type LLRB struct {
	count int
	size  int
	root  *Node
}

//-----
// tree
//-----

// New() allocates a new tree
func NewLLRB() *LLRB {
	return &LLRB{}
}

// SetRoot sets the root node of the tree.
func (t *LLRB) SetRoot(r *Node) {
	t.root = r
}

// Root returns the root node of the tree.
func (t *LLRB) Root() *Node {
	return t.root
}

// Len returns the number of nodes in the tree.
func (t *LLRB) Len() int {
	return t.count
}

// Size return the total size of keys held by this tree.
func (t *LLRB) Size() int {
	return t.size
}

//-----------------
// lookup operation
//-----------------

// Has returns true if the tree contains an element whose order
// is the same as that of key.
func (t *LLRB) Has(key Item) bool {
	return t.Get(key) != nil
}

// Get retrieves an element from the tree whose order is the
// same as that of key.
func (t *LLRB) Get(key Item) Item {
	h := t.root
	for h != nil {
		switch {
		case key.Less(h.Item):
			h = h.Left
		case h.Item.Less(key):
			h = h.Right
		default:
			return h.Item
		}
	}
	return nil
}

// Min returns the minimum element in the tree.
func (t *LLRB) Min() Item {
	h := t.root
	if h == nil {
		return nil
	}
	for h.Left != nil {
		h = h.Left
	}
	return h.Item
}

// Max returns the maximum element in the tree.
func (t *LLRB) Max() Item {
	h := t.root
	if h == nil {
		return nil
	}
	for h.Right != nil {
		h = h.Right
	}
	return h.Item
}

// UpsertBulk will upsert several keys with a single call.
// TODO: can be optimized if keys are pre-sorted.
func (t *LLRB) UpsertBulk(keys ...Item) {
	for _, key := range keys {
		t.Upsert(key)
	}
}

// InsertBulk will insert several keys with single call.
// TODO: can be optimized if keys are pre-sorted.
func (t *LLRB) InsertBulk(keys ...Item) {
	for _, key := range keys {
		t.Insert(key)
	}
}

// Upsert inserts key into the tree. If an existing
// element has the same order, it is removed from the
// tree and returned.
func (t *LLRB) Upsert(key Item) Item {
	if key == nil {
		panic("upserting nil key")
	}
	var replaced Item
	t.root, replaced = t.upsert(t.root, key)
	t.root.Black = true
	if replaced == nil {
		t.count++
	}
	return replaced
}

func (t *LLRB) upsert(h *Node, key Item) (*Node, Item) {
	if h == nil {
		return newNode(key), nil
	}

	h = walkDownRot23(h)

	var replaced Item
	if key.Less(h.Item) { // BUG
		h.Left, replaced = t.upsert(h.Left, key)
	} else if h.Item.Less(key) {
		h.Right, replaced = t.upsert(h.Right, key)
	} else {
		replaced, h.Item = h.Item, key
	}

	h = walkUpRot23(h)

	return h, replaced
}

// Insert inserts key into the tree. If an existing
// element has the same order, both elements remain in the tree.
func (t *LLRB) Insert(key Item) {
	if key == nil {
		panic("inserting nil key")
	}
	t.root = t.insert(t.root, key)
	t.root.Black = true
	t.count++
}

func (t *LLRB) insert(h *Node, key Item) *Node {
	if h == nil {
		return newNode(key)
	}

	h = walkDownRot23(h)

	if key.Less(h.Item) {
		h.Left = t.insert(h.Left, key)
	} else {
		h.Right = t.insert(h.Right, key)
	}

	return walkUpRot23(h)
}

// Rotation driver routines for 2-3 algorithm

func walkDownRot23(h *Node) *Node { return h }

func walkUpRot23(h *Node) *Node {
	if isRed(h.Right) && !isRed(h.Left) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}

// Rotation driver routines for 2-3-4 algorithm

func walkDownRot234(h *Node) *Node {
	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}

func walkUpRot234(h *Node) *Node {
	if isRed(h.Right) && !isRed(h.Left) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	return h
}

// DeleteMin deletes the minimum element in the tree and
// returns the deleted key or nil otherwise.
func (t *LLRB) DeleteMin() Item {
	var deleted Item
	t.root, deleted = deleteMin(t.root)
	if t.root != nil {
		t.root.Black = true
	}
	if deleted != nil {
		t.count--
	}
	return deleted
}

// deleteMin code for LLRB 2-3 trees
func deleteMin(h *Node) (*Node, Item) {
	if h == nil {
		return nil, nil
	}
	if h.Left == nil {
		return nil, h.Item
	}

	if !isRed(h.Left) && !isRed(h.Left.Left) {
		h = moveRedLeft(h)
	}

	var deleted Item
	h.Left, deleted = deleteMin(h.Left)

	return fixUp(h), deleted
}

// DeleteMax deletes the maximum element in the tree and
// returns the deleted key or nil otherwise
func (t *LLRB) DeleteMax() Item {
	var deleted Item
	t.root, deleted = deleteMax(t.root)
	if t.root != nil {
		t.root.Black = true
	}
	if deleted != nil {
		t.count--
	}
	return deleted
}

func deleteMax(h *Node) (*Node, Item) {
	if h == nil {
		return nil, nil
	}
	if isRed(h.Left) {
		h = rotateRight(h)
	}
	if h.Right == nil {
		return nil, h.Item
	}
	if !isRed(h.Right) && !isRed(h.Right.Left) {
		h = moveRedRight(h)
	}
	var deleted Item
	h.Right, deleted = deleteMax(h.Right)

	return fixUp(h), deleted
}

// Delete deletes an key from the tree whose key equals key.
// The deleted key is return, otherwise nil is returned.
func (t *LLRB) Delete(key Item) Item {
	var deleted Item
	t.root, deleted = t.delete(t.root, key)
	if t.root != nil {
		t.root.Black = true
	}
	if deleted != nil {
		t.count--
	}
	return deleted
}

func (t *LLRB) delete(h *Node, key Item) (*Node, Item) {
	var deleted Item
	if h == nil {
		return nil, nil
	}
	if key.Less(h.Item) {
		if h.Left == nil { // key not present. Nothing to delete
			return h, nil
		}
		if !isRed(h.Left) && !isRed(h.Left.Left) {
			h = moveRedLeft(h)
		}
		h.Left, deleted = t.delete(h.Left, key)
	} else {
		if isRed(h.Left) {
			h = rotateRight(h)
		}
		// If @key equals @h.Item and no right children at @h
		if !h.Item.Less(key) && h.Right == nil {
			return nil, h.Item
		}
		// PETAR: Added 'h.Right != nil' below
		if h.Right != nil && !isRed(h.Right) && !isRed(h.Right.Left) {
			h = moveRedRight(h)
		}
		// If @key equals @h.Item, and (from above) 'h.Right != nil'
		if !h.Item.Less(key) {
			var subDeleted Item
			h.Right, subDeleted = deleteMin(h.Right)
			if subDeleted == nil {
				panic("logic")
			}
			deleted, h.Item = h.Item, subDeleted
		} else { // Else, @key is bigger than @h.Item
			h.Right, deleted = t.delete(h.Right, key)
		}
	}

	return fixUp(h), deleted
}

//----------------
// range operation
//----------------

// KeyIterator will be called while ranging between a
// low-key and high-key
type KeyIterator func(Item) bool

// Range from a low-key and high-key, if incl is,
//  "low"  : iterate including low-key, excluding high-key
//  "high" : iterate including high-key, excluding high-key
//  "both" : iterate including both low-key and high-key
//  "none" : iterate excluding both low-key and high-key
func (t *LLRB) Range(low, high Item, incl string, iter KeyIterator) {
	switch incl {
	case "both":
		t.rangeFromFind(t.root, low, high, iter)
	case "high":
		t.rangeAfterFind(t.root, low, high, iter)
	case "low":
		t.rangeFromTill(t.root, low, high, iter)
	default:
		t.rangeAfterTill(t.root, low, high, iter)
	}
}

// low <= (keys) <= high
func (t *LLRB) rangeFromFind(h *Node, low, high Item, iter KeyIterator) bool {
	if h == nil {
		return true
	}
	if high != nil && high.Less(h.Item) {
		return t.rangeFromFind(h.Left, low, high, iter)
	}
	if low != nil && h.Item.Less(low) {
		return t.rangeFromFind(h.Right, low, high, iter)
	}
	if !t.rangeFromFind(h.Left, low, high, iter) {
		return false
	}
	if !iter(h.Item) {
		return false
	}
	return t.rangeFromFind(h.Right, low, high, iter)
}

// low <= (keys) < high
func (t *LLRB) rangeFromTill(h *Node, low, high Item, iter KeyIterator) bool {
	if h == nil {
		return true
	}
	if high != nil && !h.Item.Less(high) {
		return t.rangeFromTill(h.Left, low, high, iter)
	}
	if low != nil && h.Item.Less(low) {
		return t.rangeFromTill(h.Right, low, high, iter)
	}
	if !t.rangeFromTill(h.Left, low, high, iter) {
		return false
	}
	if !iter(h.Item) {
		return false
	}
	return t.rangeFromTill(h.Right, low, high, iter)
}

// low < (keys) <= high
func (t *LLRB) rangeAfterFind(h *Node, low, high Item, iter KeyIterator) bool {
	if h == nil {
		return true
	}
	if high != nil && high.Less(h.Item) {
		return t.rangeAfterFind(h.Left, low, high, iter)
	}
	if low != nil && !low.Less(h.Item) {
		return t.rangeAfterFind(h.Right, low, high, iter)
	}
	if !t.rangeAfterFind(h.Left, low, high, iter) {
		return false
	}
	if !iter(h.Item) {
		return false
	}
	return t.rangeAfterFind(h.Right, low, high, iter)
}

// low < (keys) < high
func (t *LLRB) rangeAfterTill(h *Node, low, high Item, iter KeyIterator) bool {
	if h == nil {
		return true
	}
	if high != nil && !h.Item.Less(high) {
		return t.rangeAfterTill(h.Left, low, high, iter)
	}
	if low != nil && !low.Less(h.Item) {
		return t.rangeAfterTill(h.Right, low, high, iter)
	}
	if !t.rangeAfterTill(h.Left, low, high, iter) {
		return false
	}
	if !iter(h.Item) {
		return false
	}
	return t.rangeAfterTill(h.Right, low, high, iter)
}

//--------------------
// statistic operation
//--------------------

// GetHeight() returns an item in the tree with
// key @key, and it's height in the tree.
func (t *LLRB) GetHeight(key Item) (result Item, depth int) {
	return getHeight(t.root, key)
}

// HeightStats() returns the average and standard
// deviation of the height of elements in the tree.
func (t *LLRB) HeightStats() (avg, stddev float64) {
	av := &Average{}
	heightStats(t.root, 0, av)
	return av.GetAvg(), av.GetStdDev()
}

func getHeight(h *Node, key Item) (Item, int) {
	if h == nil {
		return nil, 0
	}
	if key.Less(h.Item) {
		result, depth := getHeight(h.Left, key)
		return result, depth + 1
	}
	if h.Item.Less(key) {
		result, depth := getHeight(h.Right, key)
		return result, depth + 1
	}
	return h.Item, 0
}

func heightStats(h *Node, d int, av *Average) {
	if h == nil {
		return
	}
	av.Add(float64(d))
	if h.Left != nil {
		heightStats(h.Left, d+1, av)
	}
	if h.Right != nil {
		heightStats(h.Right, d+1, av)
	}
}

//-----
// node
//-----

// Node in LLRB tree.
type Node struct {
	Item
	Left, Right *Node // Pointers to left and right child nodes
	// If set, the color of the link (incoming from the parent) is black
	Black bool
	// In the LLRB, new nodes are always red, hence the zero-value for node
}

func newNode(key Item) *Node { return &Node{Item: key} }

func isRed(h *Node) bool {
	if h == nil {
		return false
	}
	return !h.Black
}

func rotateLeft(h *Node) *Node {
	x := h.Right
	if x.Black {
		panic("rotating a black link")
	}
	h.Right = x.Left
	x.Left = h
	x.Black = h.Black
	h.Black = false
	return x
}

func rotateRight(h *Node) *Node {
	x := h.Left
	if x.Black {
		panic("rotating a black link")
	}
	h.Left = x.Right
	x.Right = h
	x.Black = h.Black
	h.Black = false
	return x
}

// REQUIRE: Left and Right children must be present
func flip(h *Node) {
	h.Black = !h.Black
	h.Left.Black = !h.Left.Black
	h.Right.Black = !h.Right.Black
}

// REQUIRE: Left and Right children must be present
func moveRedLeft(h *Node) *Node {
	flip(h)
	if isRed(h.Right.Left) {
		h.Right = rotateRight(h.Right)
		h = rotateLeft(h)
		flip(h)
	}
	return h
}

// REQUIRE: Left and Right children must be present
func moveRedRight(h *Node) *Node {
	flip(h)
	if isRed(h.Left.Left) {
		h = rotateRight(h)
		flip(h)
	}
	return h
}

func fixUp(h *Node) *Node {
	if isRed(h.Right) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}
