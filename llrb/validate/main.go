package main

import "fmt"
import "flag"
import "sort"
import "time"
import "log"
import "reflect"
import "io/ioutil"
import "encoding/json"

import "github.com/prataprc/golib/llrb"
import "github.com/prataprc/goparsec"
import "github.com/prataprc/monster"
import mcommon "github.com/prataprc/monster/common"

var options struct {
	repeat int
	seed   int
	bagdir string
	mvcc   bool
}

func argParse() {
	seed := time.Now().UTC().Second()
	flag.IntVar(&options.repeat, "repeat", 1000,
		"number of times to repeat the generator")
	flag.IntVar(&options.seed, "seed", seed,
		"seed value for generating inputs")
	flag.StringVar(&options.bagdir, "bagdir", "",
		"bagdir for monster")
	flag.BoolVar(&options.mvcc, "mvcc", false,
		"use mvcc llrb")
	flag.Parse()
}

func main() {
	argParse()
	fmt.Printf("Seed: %v\n", options.seed)
	outch := make(chan [][]interface{}, 1000)
	go generate(options.repeat, "./llrb.prod", outch)
	count := options.repeat

	d := llrb.NewDict()
	var ds llrb.Api
	if options.mvcc {
		ds = llrb.NewLLRBMVCC()
	} else {
		ds = llrb.NewLLRB()
	}
	m := make(map[string]int)
	for count > 0 {
		count--
		cmds := <-outch
		for _, cmd := range cmds {
			m = validate(d, ds, cmd, m)
		}
	}
	validateEqual(d, ds)
	printStats(m)
}

func validate(
	d *llrb.Dict,
	rb llrb.Api, cmd []interface{}, m map[string]int) map[string]int {

	var ref, val interface{}
	name := cmd[0].(string)
	switch name {
	case "get":
		ref = d.Get(&llrb.KeyInt{int64(cmd[1].(float64)), -1})
		val = rb.Get(&llrb.KeyInt{int64(cmd[1].(float64)), -1})
	case "min":
		ref = d.Min()
		val = rb.Min()
	case "max":
		ref = d.Max()
		val = rb.Max()
	case "delmin":
		ref = d.DeleteMin()
		val = rb.DeleteMin()
	case "delmax":
		ref = d.DeleteMax()
		val = rb.DeleteMax()
	case "upsert":
		now := time.Now().UnixNano()
		k := &llrb.KeyInt{int64(cmd[1].(float64)), now}
		ref = d.Upsert(k)
		val = rb.Upsert(k)
	case "insert":
		now := time.Now().UnixNano()
		k := &llrb.KeyInt{int64(cmd[1].(float64)), now}
		if rb.Get(k) == nil {
			d.Insert(k)
			rb.Insert(k)
		} else {
			cmd[0] = "upsert"
			ref = d.Upsert(k)
			val = rb.Upsert(k)
		}
	case "delete":
		ref = d.Delete(&llrb.KeyInt{int64(cmd[1].(float64)), -1})
		val = rb.Delete(&llrb.KeyInt{int64(cmd[1].(float64)), -1})
	default:
		log.Fatalf("unknown command %v\n", cmd)
		return m
	}
	if reflect.DeepEqual(ref, val) == false {
		log.Fatalf("expected %v got %v\n", ref, val)
	}
	m = stats(cmd, m)
	return m
}

func validateEqual(d *llrb.Dict, rb llrb.Api) {
	refKeys := make([]*llrb.KeyInt, 0)
	fmt.Printf("number of elements {%v,%v}\n", d.Len(), rb.Len())
	rb.Range(nil, nil, "both", func(k llrb.Item) bool {
		refKeys = append(refKeys, k.(*llrb.KeyInt))
		return true
	})
	keys := make([]*llrb.KeyInt, 0)
	d.Range(nil, nil, "both", func(k llrb.Item) bool {
		keys = append(keys, k.(*llrb.KeyInt))
		return true
	})
	if reflect.DeepEqual(refKeys, keys) == false {
		log.Fatalf("final Dict keys and LLRB keys mismatch\n")
	}
}

func generate(repeat int, prodfile string, outch chan<- [][]interface{}) {
	text, err := ioutil.ReadFile(prodfile)
	if err != nil {
		log.Fatal(err)
	}
	root := compile(parsec.NewScanner(text)).(mcommon.Scope)
	seed, bagdir, prodfile := uint64(options.seed), options.bagdir, prodfile
	scope := monster.BuildContext(root, seed, bagdir, prodfile)
	nterms := scope["_nonterminals"].(mcommon.NTForms)
	for i := 0; i < repeat; i++ {
		scope = scope.RebuildContext()
		val := evaluate("root", scope, nterms["s"])
		var arr [][]interface{}
		if err := json.Unmarshal([]byte(val.(string)), &arr); err != nil {
			log.Fatal(err)
		}
		outch <- arr
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

func stats(cmd []interface{}, m map[string]int) map[string]int {
	name := cmd[0].(string)
	count, ok := m[name]
	if !ok {
		count = 0
	}
	m[name] = count + 1
	return m
}

func printStats(m map[string]int) {
	keys, total := []string{}, 0
	for name, count := range m {
		keys = append(keys, name)
		total += count
	}
	fmt.Printf("total commands : %v\n", total)
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%v: %v\n", key, m[key])
	}
}
