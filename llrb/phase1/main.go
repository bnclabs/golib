package main

import "fmt"
import "flag"
import "time"

//import "math/rand"

//import "github.com/prataprc/golib/llrb"
import "github.com/prataprc/goparsec"
import "github.com/prataprc/monster"
import mcommon "github.com/prataprc/monster/common"

var options struct {
	byt     bool
	num     bool
	seed    int
	entries int
	bagdir  string
}

func argParse() {
	seed := time.Now().UTC().Second()
	flag.BoolVar(&options.byt, "bytes", false,
		"use bytes as keys")
	flag.BoolVar(&options.num, "num", true,
		"use numbers as keys")
	flag.IntVar(&options.entries, "entries", 1000,
		"number of entries to work with")
	flag.IntVar(&options.seed, "seed", seed,
		"seed value for generating inputs")
	flag.StringVar(&options.bagdir, "bagdir", "",
		"bagdir for monster")
	flag.Parse()
}

func main() {
	argParse()
	took := 0
	if options.num {
		fmt.Printf("took %v to generate %v entries\n", took, options.entries)
	}
}

func generateString(text []byte, count int, prodfile string, outch chan<- []byte) {
	root := compile(parsec.NewScanner(text)).(mcommon.Scope)
	seed, bagdir, prodfile := uint64(options.seed), options.bagdir, prodfile
	scope := monster.BuildContext(root, seed, bagdir, prodfile)
	nterms := scope["_nonterminals"].(mcommon.NTForms)
	for i := 0; i < count; i++ {
		scope = scope.RebuildContext()
		val := evaluate("root", scope, nterms["s"])
		outch <- []byte(val.(string))
	}
}

func compile(s parsec.Scanner) parsec.ParsecNode {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%v at %v", r, s.GetCursor())
		}
	}()
	root, _ := monster.Y(s)
	return root
}

func evaluate(
	name string, scope mcommon.Scope, forms []*mcommon.Form) interface{} {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%v", r)
		}
	}()
	return monster.EvalForms(name, scope, forms)
}

func timeit(fn func()) float64 {
	start := time.Now()
	fn()
	return float64(time.Since(start).Nanoseconds()) / 1000000
}
