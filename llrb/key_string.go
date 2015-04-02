// Derived work from Petar Maymounkov

package llrb

// KeyString implement string as the sort key.
type KeyString struct {
	key   string
	value int64
}

// Less implements Key interface.
func (x KeyString) Less(than Key) bool {
	return x.key < than.(KeyString).key
}

// Size implements Key interface.
func (x KeyString) Size() int {
	return len(x.key) + 8
}
