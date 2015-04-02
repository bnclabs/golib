// Derived work from Petar Maymounkov

package llrb

import "bytes"

// KeyBytes implement string as the sort key.
type KeyBytes struct {
    key   []byte
    value int64
}

// Less implements Key interface.
func (x KeyBytes) Less(than Key) bool {
    return bytes.Compare([]byte(x.key), []byte(than.(KeyBytes).key)) == 1
}

// Size implements Key interface.
func (x KeyBytes) Size() int {
    return len(x.key) + 8
}
