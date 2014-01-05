package lists

type Generic interface{}

func Map(fn func(Generic) Generic, xs []Generic) []Generic {
	result := make([]Generic, len(xs))
	for i, x := range xs {
		result[i] = fn(x)
	}
	return result
}

func Reduce(fn func(Generic, Generic) Generic, xs []Generic, a Generic) Generic {
	for _, x := range xs {
		a = fn(x, a)
	}
	return a
}

func Filter(fn func(Generic) bool, xs []Generic) []Generic {
	result := make([]Generic, len(xs))
	var truth, i, j = false, 0, 0
	// TODO : if fn is nil we must provide a default function that check for
	for ; i < len(xs); i = i + 1 {
		if fn == nil {
			truth = fn(xs[i])
		} else {
			truth = checkZV(xs[i])
		}
		if truth {
			result[j] = xs[i]
			j++
		}
	}
	return result[:j]
}

func checkZV(x Generic) bool {
	// TODO : make this more generic
	var y = x.(int)
	return y != 0
}

func Sum(xs []Generic) int64 {
	fn := func(x Generic, a Generic) Generic { return a.(int) + x.(int) }
	return Reduce(fn, xs, 0).(int64)
}

func Seq(from int, to int, vars ...interface{}) []int {
	result := make([]int, to-from+1)
	step := 1
	j := 0
	if len(vars) > 0 {
		step = vars[0].(int)
	}
	for i := from; i < to; i, j = i+step, j+1 {
		result[j] = i
	}
	return result[:j]
}

func Dot(f func(x Generic) Generic, g func(x Generic) Generic) func(Generic) Generic {
	return func(x Generic) Generic { return f(g(x)) }
}
