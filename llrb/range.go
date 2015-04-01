package llrb

// KeyIterator will be called while ranging between a
// low-key and high-key
type KeyIterator func(Key) bool

// Range from a low-key and high-key, if incl is,
//  "low"  : iterate including low-key, excluding high-key
//  "high" : iterate including high-key, excluding high-key
//  "both" : iterate including both low-key and high-key
//  "none" : iterate excluding both low-key and high-key
func (t *LLRB) Range(low, high Key, incl string, iter KeyIterator) {
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
func (t *LLRB) rangeFromFind(h *Node, low, high Key, iter KeyIterator) bool {
	if h == nil {
		return true
	}
	if high.Less(h.Key) {
		return t.rangeFromFind(h.Left, low, high, iter)
	}
	if h.Key.Less(low) {
		return t.rangeFromFind(h.Right, low, high, iter)
	}
	if !t.rangeFromFind(h.Left, low, high, iter) {
		return false
	}
	if !iter(h.Key) {
		return false
	}
	return t.rangeFromFind(h.Right, low, high, iter)
}

// low <= (keys) < high
func (t *LLRB) rangeFromTill(h *Node, low, high Key, iter KeyIterator) bool {
	if h == nil {
		return true
	}
	if !h.Key.Less(high) {
		return t.rangeFromTill(h.Left, low, high, iter)
	}
	if h.Key.Less(low) {
		return t.rangeFromTill(h.Right, low, high, iter)
	}
	if !t.rangeFromTill(h.Left, low, high, iter) {
		return false
	}
	if !iter(h.Key) {
		return false
	}
	return t.rangeFromTill(h.Right, low, high, iter)
}

// low < (keys) <= high
func (t *LLRB) rangeAfterFind(h *Node, low, high Key, iter KeyIterator) bool {
	if h == nil {
		return true
	}
	if high.Less(h.Key) {
		return t.rangeAfterFind(h.Left, low, high, iter)
	}
	if !low.Less(h.Key) {
		return t.rangeAfterFind(h.Right, low, high, iter)
	}
	if !t.rangeAfterFind(h.Left, low, high, iter) {
		return false
	}
	if !iter(h.Key) {
		return false
	}
	return t.rangeAfterFind(h.Right, low, high, iter)
}

// low < (keys) < high
func (t *LLRB) rangeAfterTill(h *Node, low, high Key, iter KeyIterator) bool {
	if h == nil {
		return true
	}
	if !h.Key.Less(high) {
		return t.rangeAfterTill(h.Left, low, high, iter)
	}
	if !low.Less(h.Key) {
		return t.rangeAfterTill(h.Right, low, high, iter)
	}
	if !t.rangeAfterTill(h.Left, low, high, iter) {
		return false
	}
	if !iter(h.Key) {
		return false
	}
	return t.rangeAfterTill(h.Right, low, high, iter)
}
