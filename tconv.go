package golib

type Interface interface{}

func Int(val Interface, defval int) int {
    if ival, ok := val.(int); ok {
        return ival
    } else {
        return defval
    }
}

func Bool(val Interface, defval bool) bool {
    if bval, ok := val.(bool); ok {
        return bval
    } else {
        return defval
    }
}

func String(val Interface, defval string) string {
    if sval, ok := val.(string); ok {
        return sval
    } else {
        return defval
    }
}
