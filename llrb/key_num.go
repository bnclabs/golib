// Derived work from Petar Maymounkov

package llrb

// KeyInt implements int64 as the sort key.
type KeyInt struct {
	key   int64
	value int64
}

// Less implements Key interface.
func (x KeyInt) Less(than Key) bool {
	return x.key < than.(KeyInt).key
}

// Size implements Key interface.
func (x KeyInt) Size() int {
	return 16
}
