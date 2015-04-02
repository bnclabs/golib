package llrb

// limits of a key space is defined from nInf to pInf
var ninf = nInf{}
var pinf = pInf{}

type nInf struct{}

// Less implements the Item interface.
func (nInf) Less(Item) bool {
	return true
}

// Size implements the Item interface.
func (nInf) Size() int {
	return 0
}

type pInf struct{}

// Less implements the key interface.
func (pInf) Less(Item) bool {
	return false
}

// Size implements the key interface.
func (pInf) Size() int {
	return 0
}

// Inf returns an Item that is "bigger than" any other key,
// if sign is positive. Otherwise  it returns an Item that
// is "smaller than" any other key.
func Inf(sign int) Item {
	if sign == 0 {
		panic("sign")
	}
	if sign > 0 {
		return pinf
	}
	return ninf
}
