package llrb

// MemStore exposed by in-memory data-structure implementing
// a sorted key-value store.
type MemStore interface {
	// Len return number of entries.
	Len() int

	// Has return whether Item is present in the store.
	Has(key Item) bool

	// Get return the Key-Value entry.
	Get(key Item) Item

	// Min return entry with lowest order.
	Min() Item

	// Max return entry with highest order.
	Max() Item

	// UpsertBulk will upsert 1 or more Key-Value entries.
	UpsertBulk(keys ...Item)

	// InsertBulk will insert 1 or more Key-Value entries.
	InsertBulk(keys ...Item)

	// Upsert will upsert 1 Key-Value entry.
	Upsert(key Item) Item

	// Insert will insert 1 Key-Value entry.
	Insert(key Item)

	// DeleteMin will remove Key-Value entry with lowest order.
	DeleteMin() Item

	// DeleteMax will remove Key-Value entry with highest order.
	DeleteMax() Item

	// Delete will remove the Key-Value entry with specified order.
	Delete(key Item) Item

	// Range will return a subset of sorted Key-Value entries.
	Range(low, high Item, incl string, iter KeyIterator)

	// GetHeight return the depth of a Key-Value entry with
	// specified order.
	GetHeight(key Item) (result Item, depth int)

	// HeightStats return the average depth and standard-deviation
	// of depths of all Key-Value entries.
	HeightStats() (avg, stddev float64)

	// RSnapshot shall be called on writer instance and return
	// a new snapshot instance that won't be disturbed by future
	// writes.
	RSnapshot(timeout int) MemStore

	// Release this snapshot and its resources
	ReleaseSnapshot()
}

// Item implements an Key-Value entry in the sorted list.
type Item interface {
	// Size returns the size f data held by this object.
	Size() int

	// Less return true of this key lesser.
	Less(than Item) bool
}
