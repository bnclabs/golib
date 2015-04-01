package llrb

// limits of a key space is defined from nInf to pInf
var ninf = nInf{}
var pinf = pInf{}

type nInf struct{}

// Less implements the Key interface.
func (nInf) Less(Key) bool {
	return true
}

// Size implements the Key interface.
func (nInf) Size() int {
	return 0
}

type pInf struct{}

// Less implements the key interface.
func (pInf) Less(Key) bool {
	return false
}

// Size implements the key interface.
func (pInf) Size() int {
	return 0
}

// Inf returns an Key that is "bigger than" any other key,
// if sign is positive. Otherwise  it returns an Key that
// is "smaller than" any other key.
func Inf(sign int) Key {
	if sign == 0 {
		panic("sign")
	}
	if sign > 0 {
		return pinf
	}
	return ninf
}
