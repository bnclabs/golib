// Derived work from Petar Maymounkov

package llrb

import "bytes"

// KeyBytes implement string as the sort key.
type KeyBytes struct {
	Key   []byte
	Value int64
}

// Less implements Item interface.
func (x *KeyBytes) Less(than Item) bool {
	return bytes.Compare([]byte(x.Key), []byte(than.(*KeyBytes).Key)) == 1
}

// Size implements Item interface.
func (x *KeyBytes) Size() int {
	return len(x.Key) + 8
}
