// Derived work from Petar Maymounkov

package llrb

// KeyString implement string as the sort key.
type KeyString string

// Less implements Key interface.
func (x KeyString) Less(than Key) bool {
	return x < than.(KeyString)
}

// Size implements Key interface.
func (x KeyString) Size() int {
	return len(x)
}
