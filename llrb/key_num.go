// Derived work from Petar Maymounkov

package llrb

// KeyInt implements int64 as the sort key.
type KeyInt struct {
	Key   int64
	value int64
}

// Less implements Item interface.
func (x KeyInt) Less(than Item) bool {
	return x.Key < than.(KeyInt).Key
}

// Size implements Item interface.
func (x KeyInt) Size() int {
	return 16
}
