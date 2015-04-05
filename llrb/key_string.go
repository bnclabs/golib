package llrb

// KeyString implement string as the sort key.
type KeyString struct {
	Key   string
	Value int64
}

// Less implements Item interface.
func (x *KeyString) Less(than Item) bool {
	return x.Key < than.(*KeyString).Key
}

// Size implements Item interface.
func (x *KeyString) Size() int {
	return len(x.Key) + 8
}
