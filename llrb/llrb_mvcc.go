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

import "fmt"

var _ = fmt.Sprintf("dummy format")

// LLRBMVCC is a Left-Leaning Red-Black (LLRB)
// implementation of 2-3 trees supporting concurrent
// reads.
type LLRBMVCC struct {
	count int
	size  int
	root  *Node
}

//-----
// tree
//-----

// New() allocates a new tree.
func NewLLRBMVCC() *LLRBMVCC {
	return &LLRBMVCC{}
}

// SetRoot sets the root node of the tree.
func (t *LLRBMVCC) SetRoot(r *Node) {
	t.root = r
}

// Root returns the root node of the tree.
func (t *LLRBMVCC) Root() *Node {
	return t.root
}

// Len returns the number of nodes in the tree.
func (t *LLRBMVCC) Len() int {
	return t.count
}

// Size return the total size of keys held by this tree.
func (t *LLRBMVCC) Size() int {
	return t.size
}

//-----------------
// lookup operation
//-----------------

// Has returns true if the tree contains an element whose order
// is the same as that of key.
func (t *LLRBMVCC) Has(key Item) bool {
	return t.Get(key) != nil
}

// Get retrieves an element from the tree whose order is the
// same as that of key.
func (t *LLRBMVCC) Get(key Item) Item {
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
func (t *LLRBMVCC) Min() Item {
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
func (t *LLRBMVCC) Max() Item {
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
func (t *LLRBMVCC) UpsertBulk(keys ...Item) {
	for _, key := range keys {
		t.Upsert(key)
	}
}

// InsertBulk will insert several keys with single call.
// TODO: can be optimized if keys are pre-sorted.
func (t *LLRBMVCC) InsertBulk(keys ...Item) {
	for _, key := range keys {
		t.Insert(key)
	}
}

// Upsert inserts key into the tree. If an existing
// element has the same order, it is removed from the
// tree and returned.
func (t *LLRBMVCC) Upsert(key Item) Item {
	if key == nil {
		panic("upserting nil key")
	}
	var replaced Item
	reclaim := []*Node{}
	t.root, replaced, reclaim = t.upsert(t.root, key, reclaim)
	t.root.Black = true
	if replaced == nil {
		t.count++
	}
	return replaced
}

func (t *LLRBMVCC) upsert(
	h *Node, key Item, reclaim []*Node) (*Node, Item, []*Node) {

	if h == nil {
		return newNode(key), nil, reclaim
	}

	reclaim = append(reclaim, h)
	hnew := cow(h)

	hnew = walkDownRot23COW(hnew)

	var replaced Item
	if key.Less(hnew.Item) { // BUG
		hnew.Left, replaced, reclaim = t.upsert(hnew.Left, key, reclaim)
	} else if hnew.Item.Less(key) {
		hnew.Right, replaced, reclaim = t.upsert(hnew.Right, key, reclaim)
	} else {
		replaced, hnew.Item = hnew.Item, key
	}

	hnew, reclaim = walkUpRot23COW(hnew, reclaim)
	return hnew, replaced, reclaim
}

// Insert inserts key into the tree. If an existing
// element has the same order, both elements remain in the tree.
func (t *LLRBMVCC) Insert(key Item) {
	if key == nil {
		panic("inserting nil key")
	}
	reclaim := []*Node{}
	t.root, reclaim = t.insert(t.root, key, reclaim)
	t.root.Black = true
	t.count++
}

func (t *LLRBMVCC) insert(h *Node, key Item, reclaim []*Node) (*Node, []*Node) {
	if h == nil {
		return newNode(key), reclaim
	}

	reclaim = append(reclaim, h)
	hnew := cow(h)

	hnew = walkDownRot23COW(hnew)

	if key.Less(hnew.Item) {
		hnew.Left, reclaim = t.insert(hnew.Left, key, reclaim)
	} else {
		hnew.Right, reclaim = t.insert(hnew.Right, key, reclaim)
	}

	hnew, reclaim = walkUpRot23COW(hnew, reclaim)
	return hnew, reclaim
}

// Rotation driver routines for 2-3 algorithm

func walkDownRot23COW(hnew *Node) *Node { return hnew }

func walkUpRot23COW(hnew *Node, reclaim []*Node) (*Node, []*Node) {
	if isRed(hnew.Right) && !isRed(hnew.Left) {
		hnew, reclaim = rotateLeftCOW(hnew, reclaim)
	}

	if isRed(hnew.Left) && isRed(hnew.Left.Left) {
		hnew, reclaim = rotateRightCOW(hnew, reclaim)
	}

	if isRed(hnew.Left) && isRed(hnew.Right) {
		reclaim = flipCOW(hnew, reclaim)
	}

	return hnew, reclaim
}

// DeleteMin deletes the minimum element in the tree and
// returns the deleted key or nil otherwise.
func (t *LLRBMVCC) DeleteMin() Item {
	var deleted Item
	reclaim := []*Node{}
	t.root, deleted, reclaim = deleteMinCOW(t.root, reclaim)
	if t.root != nil {
		t.root.Black = true
	}
	if deleted != nil {
		t.count--
	}
	return deleted
}

// deleteMinCOW code for LLRBMVCC 2-3 trees
func deleteMinCOW(h *Node, reclaim []*Node) (*Node, Item, []*Node) {
	if h == nil {
		return nil, nil, reclaim
	}
	if h.Left == nil {
		reclaim = append(reclaim, h)
		return nil, h.Item, reclaim
	}

	reclaim = append(reclaim, h)
	hnew := cow(h)

	if !isRed(hnew.Left) && !isRed(hnew.Left.Left) {
		hnew, reclaim = moveRedLeftCOW(hnew, reclaim)
	}

	var deleted Item
	hnew.Left, deleted, reclaim = deleteMinCOW(hnew.Left, reclaim)

	hnew, reclaim = fixUpCOW(hnew, reclaim)
	return hnew, deleted, reclaim
}

// DeleteMax deletes the maximum element in the tree and
// returns the deleted key or nil otherwise
func (t *LLRBMVCC) DeleteMax() Item {
	var deleted Item
	reclaim := []*Node{}
	t.root, deleted, reclaim = deleteMaxCOW(t.root, reclaim)
	if t.root != nil {
		t.root.Black = true
	}
	if deleted != nil {
		t.count--
	}
	return deleted
}

func deleteMaxCOW(h *Node, reclaim []*Node) (*Node, Item, []*Node) {
	if h == nil {
		return nil, nil, reclaim
	}

	reclaim = append(reclaim, h)
	hnew := cow(h)

	if isRed(hnew.Left) {
		hnew, reclaim = rotateRightCOW(hnew, reclaim)
	}
	if hnew.Right == nil {
		return nil, hnew.Item, reclaim
	}
	if !isRed(hnew.Right) && !isRed(hnew.Right.Left) {
		hnew, reclaim = moveRedRightCOW(hnew, reclaim)
	}
	var deleted Item
	hnew.Right, deleted, reclaim = deleteMaxCOW(hnew.Right, reclaim)

	hnew, reclaim = fixUpCOW(hnew, reclaim)
	return hnew, deleted, reclaim
}

// Delete deletes an key from the tree whose key equals key.
// The deleted key is return, otherwise nil is returned.
func (t *LLRBMVCC) Delete(key Item) Item {
	var deleted Item
	reclaim := []*Node{}
	t.root, deleted, reclaim = t.delete(t.root, key, reclaim)
	if t.root != nil {
		t.root.Black = true
	}
	if deleted != nil {
		t.count--
	}
	return deleted
}

func (t *LLRBMVCC) delete(
	h *Node, key Item, reclaim []*Node) (*Node, Item, []*Node) {

	var deleted Item
	if h == nil {
		return nil, nil, reclaim
	}

	reclaim = append(reclaim, h)
	hnew := cow(h)

	if key.Less(hnew.Item) {
		if hnew.Left == nil { // key not present. Nothing to delete
			return hnew, nil, reclaim
		}
		if !isRed(hnew.Left) && !isRed(hnew.Left.Left) {
			hnew, reclaim = moveRedLeftCOW(hnew, reclaim)
		}
		hnew.Left, deleted, reclaim = t.delete(hnew.Left, key, reclaim)
	} else {
		if isRed(hnew.Left) {
			hnew, reclaim = rotateRightCOW(hnew, reclaim)
		}
		// If @key equals @hnew.Item and no right children at @hnew
		if !hnew.Item.Less(key) && hnew.Right == nil {
			return nil, hnew.Item, reclaim
		}
		// PETAR: Added 'hnew.Right != nil' below
		if hnew.Right != nil && !isRed(hnew.Right) && !isRed(hnew.Right.Left) {
			hnew, reclaim = moveRedRightCOW(hnew, reclaim)
		}
		// If @key equals @hnew.Item, and (from above) 'hnew.Right != nil'
		if !hnew.Item.Less(key) {
			var subDeleted Item
			hnew.Right, subDeleted, reclaim = deleteMinCOW(hnew.Right, reclaim)
			if subDeleted == nil {
				panic("logic")
			}
			deleted, hnew.Item = hnew.Item, subDeleted
		} else { // Else, @key is bigger than @hnew.Item
			hnew.Right, deleted, reclaim = t.delete(hnew.Right, key, reclaim)
		}
	}

	hnew, reclaim = fixUpCOW(hnew, reclaim)
	return hnew, deleted, reclaim
}

//----------------
// range operation
//----------------

// Range from a low-key and high-key, if incl is,
//  "low"  : iterate including low-key, excluding high-key
//  "high" : iterate including high-key, excluding high-key
//  "both" : iterate including both low-key and high-key
//  "none" : iterate excluding both low-key and high-key
func (t *LLRBMVCC) Range(low, high Item, incl string, iter KeyIterator) {
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
func (t *LLRBMVCC) rangeFromFind(h *Node, low, high Item, iter KeyIterator) bool {
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
func (t *LLRBMVCC) rangeFromTill(h *Node, low, high Item, iter KeyIterator) bool {
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
func (t *LLRBMVCC) rangeAfterFind(h *Node, low, high Item, iter KeyIterator) bool {
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
func (t *LLRBMVCC) rangeAfterTill(h *Node, low, high Item, iter KeyIterator) bool {
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

// GetHeight() returns an item in the tree with key @key,
// and it's height in the tree
func (t *LLRBMVCC) GetHeight(key Item) (result Item, depth int) {
	return t.getHeight(t.root, key)
}

// HeightStats() returns the average and standard deviation of the height
// of elements in the tree
func (t *LLRBMVCC) HeightStats() (avg, stddev float64) {
	av := &avgVar{}
	heightStats(t.root, 0, av)
	return av.GetAvg(), av.GetStdDev()
}

//------------------------
// copy on write operation
//------------------------

func cow(h *Node) *Node {
	if h == nil {
		return h
	}
	hnew := &Node{
		Item:  h.Item,
		Left:  h.Left,
		Right: h.Right,
		Black: h.Black,
	}
	return hnew
}

func rotateLeftCOW(hnew *Node, reclaim []*Node) (*Node, []*Node) {
	reclaim = append(reclaim, hnew.Right)
	y := cow(hnew.Right)
	if y.Black {
		panic("rotating a black link")
	}
	hnew.Right = y.Left
	y.Left = hnew
	y.Black = hnew.Black
	hnew.Black = false
	return y, reclaim
}

func rotateRightCOW(hnew *Node, reclaim []*Node) (*Node, []*Node) {
	reclaim = append(reclaim, hnew.Left)
	x := cow(hnew.Left)
	if x.Black {
		panic("rotating a black link")
	}
	hnew.Left = x.Right
	x.Right = hnew
	x.Black = hnew.Black
	hnew.Black = false
	return x, reclaim
}

// REQUIRE: Left and Right children must be present
func flipCOW(hnew *Node, reclaim []*Node) []*Node {
	reclaim = append(reclaim, hnew.Left, hnew.Right)
	x, y := cow(hnew.Left), cow(hnew.Right)
	x.Black = !x.Black
	y.Black = !y.Black
	hnew.Black, hnew.Left, hnew.Right = !hnew.Black, x, y
	return reclaim
}

// REQUIRE: Left and Right children must be present
func moveRedLeftCOW(hnew *Node, reclaim []*Node) (*Node, []*Node) {
	reclaim = flipCOW(hnew, reclaim)
	if isRed(hnew.Right.Left) {
		hnew.Right, reclaim = rotateRightCOW(hnew.Right, reclaim)
		hnew, reclaim = rotateLeftCOW(hnew, reclaim)
		reclaim = flipCOW(hnew, reclaim)
	}
	return hnew, reclaim
}

// REQUIRE: Left and Right children must be present
func moveRedRightCOW(hnew *Node, reclaim []*Node) (*Node, []*Node) {
	reclaim = flipCOW(hnew, reclaim)
	if isRed(hnew.Left.Left) {
		hnew, reclaim = rotateRightCOW(hnew, reclaim)
		reclaim = flipCOW(hnew, reclaim)
	}
	return hnew, reclaim
}

func fixUpCOW(hnew *Node, reclaim []*Node) (*Node, []*Node) {
	if isRed(hnew.Right) {
		hnew, reclaim = rotateLeftCOW(hnew, reclaim)
	}

	if isRed(hnew.Left) && isRed(hnew.Left.Left) {
		hnew, reclaim = rotateRightCOW(hnew, reclaim)
	}

	if isRed(hnew.Left) && isRed(hnew.Right) {
		reclaim = flipCOW(hnew, reclaim)
	}

	return hnew, reclaim
}
