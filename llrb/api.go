package llrb

// Api exposed by in-memory data-structure implementing
// a sorted key-value store.
type Api interface {
	Len() int
	Has(key Item) bool
	Get(key Item) Item
	Min() Item
	Max() Item
	UpsertBulk(keys ...Item)
	InsertBulk(keys ...Item)
	Upsert(key Item) Item
	Insert(key Item)
	DeleteMin() Item
	DeleteMax() Item
	Delete(key Item) Item
	Range(low, high Item, incl string, iter KeyIterator)
	GetHeight(key Item) (result Item, depth int)
	getHeight(h *Node, key Item) (Item, int)
	HeightStats() (avg, stddev float64)
}

// Item implements an Key-Value entry in the sorted list.
type Item interface {
	// Size returns the size f data held by this object.
	Size() int

	// Less return true of this key lesser.
	Less(than Item) bool
}
