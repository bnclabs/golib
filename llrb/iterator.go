package llrb

func (t *LLRB) AscendRange(greaterOrEqual, lessThan Key, iterator KeyIterator) {
	t.ascendRange(t.root, greaterOrEqual, lessThan, iterator)
}

func (t *LLRB) ascendRange(h *Node, inf, sup Key, iterator KeyIterator) bool {
	if h == nil {
		return true
	}
	if !h.Key.Less(sup) {
		return t.ascendRange(h.Left, inf, sup, iterator)
	}
	if h.Key.Less(inf) {
		return t.ascendRange(h.Right, inf, sup, iterator)
	}

	if !t.ascendRange(h.Left, inf, sup, iterator) {
		return false
	}
	if !iterator(h.Key) {
		return false
	}
	return t.ascendRange(h.Right, inf, sup, iterator)
}

// AscendGreaterOrEqual will call iterator once for each element greater
// or equal to pivot in ascending order. It will stop whenever the iterator
// returns false.
func (t *LLRB) AscendGreaterOrEqual(pivot Key, iterator KeyIterator) {
	t.ascendGreaterOrEqual(t.root, pivot, iterator)
}

func (t *LLRB) ascendGreaterOrEqual(h *Node, pivot Key, iterator KeyIterator) bool {
	if h == nil {
		return true
	}
	if !h.Key.Less(pivot) {
		if !t.ascendGreaterOrEqual(h.Left, pivot, iterator) {
			return false
		}
		if !iterator(h.Key) {
			return false
		}
	}
	return t.ascendGreaterOrEqual(h.Right, pivot, iterator)
}

func (t *LLRB) AscendLessThan(pivot Key, iterator KeyIterator) {
	t.ascendLessThan(t.root, pivot, iterator)
}

func (t *LLRB) ascendLessThan(h *Node, pivot Key, iterator KeyIterator) bool {
	if h == nil {
		return true
	}
	if !t.ascendLessThan(h.Left, pivot, iterator) {
		return false
	}
	if !iterator(h.Key) {
		return false
	}
	if h.Key.Less(pivot) {
		return t.ascendLessThan(h.Left, pivot, iterator)
	}
	return true
}

// DescendLessOrEqual will call iterator once for each element less than the
// pivot in descending order. It will stop whenever the iterator returns false.
func (t *LLRB) DescendLessOrEqual(pivot Key, iterator KeyIterator) {
	t.descendLessOrEqual(t.root, pivot, iterator)
}

func (t *LLRB) descendLessOrEqual(h *Node, pivot Key, iterator KeyIterator) bool {
	if h == nil {
		return true
	}
	if h.Key.Less(pivot) || !pivot.Less(h.Key) {
		if !t.descendLessOrEqual(h.Right, pivot, iterator) {
			return false
		}
		if !iterator(h.Key) {
			return false
		}
	}
	return t.descendLessOrEqual(h.Left, pivot, iterator)
}
