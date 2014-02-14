package golib

// Take an interface `val`, that points to `int` type and return the integer
// value. Otherwise return `defval`
func Int(val Interface, defval int) int {
	if ival, ok := val.(int); ok {
		return ival
	} else {
		return defval
	}
}

// Take an interface `val`, that points to `bool` type and return the boolean
// value. Otherwise return `defval`
func Bool(val Interface, defval bool) bool {
	if bval, ok := val.(bool); ok {
		return bval
	} else {
		return defval
	}
}

// Take an interface `val`, that points to `string` type and return the string
// value. Otherwise return `defval`
func String(val Interface, defval string) string {
	if sval, ok := val.(string); ok {
		return sval
	}
	return defval
}
