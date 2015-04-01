// Derived work from Petar Maymounkov

package llrb

// KeyInt implements int64 as the sort key.
type KeyInt int64

// Less implements Key interface.
func (x KeyInt) Less(than Key) bool {
	return x < than.(KeyInt)
}

// Size implements Key interface.
func (x KeyInt) Size() int {
	return 8
}
