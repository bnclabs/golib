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
	flag.StringVar(&options.bagdir, "bagdir", "./",
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

	if options.mvcc {
		withLLRBMVCC(options.repeat, outch)
	} else {
		withLLRB(options.repeat, outch)
	}
}

func withLLRB(count int, outch chan [][]interface{}) {
	d := llrb.NewDict()
	ds := llrb.NewLLRB()
	stats := make(map[string]int)
	for count > 0 {
		count--
		cmds := <-outch
		for _, cmd := range cmds {
			stats = validate(d, ds, cmd, stats)
		}
	}
	validateEqual(d, ds)
	printStats(stats)
	avg, sd := ds.HeightStats()
	fmt.Printf("LLRB Stats: avg-height: %4.2f, sd-height: %4.2f\n", avg, sd)
}

func withLLRBMVCC(count int, outch chan [][]interface{}) {
	d := llrb.NewDict()
	w := llrb.NewLLRBMVCC(10)
	quitch := make(chan bool)
	readers := make(map[int][]interface{})
	for i := 0; i < 4; i++ {
		inpch := make(chan []interface{}, 4)
		rstats := make(map[string]int)
		readers[i] = []interface{}{inpch, rstats}
		go concurrent_reader(inpch, rstats, quitch)
		inpch <- []interface{}{"snapshot", d.RSnapshot(100), w.RSnapshot(100)}
	}

	total := 0
	stats := make(map[string]int)
	snapstick := time.Tick(100 * time.Millisecond) // take snapshot per 100ms
	for count > 0 {
		count--
		select {
		case cmds := <-outch:
			for _, cmd := range cmds {
				total++
				if isReadOp(cmd) {
					for _, reader := range readers {
						reader[0].(chan []interface{}) <- cmd
					}
				} else {
					stats = validate(d, w, cmd, stats)
				}
			}

		case <-snapstick:
			for i, reader := range readers {
				if i >= total%4 {
					break
				}
				d := d.RSnapshot(100)
				ds := w.RSnapshot(100)
				reader[0].(chan []interface{}) <- []interface{}{"snapshot", d, ds}
			}
		}
	}

	// close readers
	for _, reader := range readers {
		close(reader[0].(chan []interface{}))
	}
	// wait for reader routines to quit.
	count = len(readers)
	for count > 0 {
		<-quitch
		count--
	}

	validateEqual(d, w)
	fmt.Printf("total number of commands: %v\n", total)
	fmt.Println("stats for writer:")
	printStats(stats)
	for _, reader := range readers {
		fmt.Println("stats for reader:")
		printStats(reader[1].(map[string]int))
	}
	avg, sd := w.HeightStats()
	fmt.Printf("LLRB Stats: avg-height: %4.2f, sd-height: %4.2f\n", avg, sd)
	for opname := range writeOps {
		avg, sd := w.CowStats(opname)
		fmt.Printf("COW %s: avg-cow: %4.2f, sd-cow: %4.2f\n", opname, avg, sd)
	}
}

func concurrent_reader(
	inpch chan []interface{}, stats map[string]int, quitch chan bool) {

	d := (*llrb.Dict)(nil)
	ds := (*llrb.LLRBMVCC)(nil)
	for {
		cmd, ok := <-inpch
		if !ok {
			break
		}
		name := cmd[0].(string)
		if name == "snapshot" {
			if ds != nil {
				ds.ReleaseSnapshot()
			}
			d = cmd[1].(*llrb.Dict)
			ds = cmd[2].(*llrb.LLRBMVCC)
			stats["snapshot"]++
		} else if isReadOp(cmd) {
			if ds == nil || d == nil {
				log.Fatalf("reader not initialized with snapshot")
			}
			validate(d, ds, cmd, stats)

		} else {
			log.Fatalf("write op not allowed")
		}
	}
	quitch <- true
}

//--------
// monster
//--------

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

//---------
// validate
//---------

func validate(
	d *llrb.Dict,
	rb llrb.MemStore, cmd []interface{}, stats map[string]int) map[string]int {

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
		return stats
	}
	if reflect.DeepEqual(ref, val) == false {
		log.Fatalf("expected %v got %v\n", ref, val)
	}
	stats = compute_stats(cmd, stats)
	return stats
}

func validateEqual(d *llrb.Dict, rb llrb.MemStore) {
	refKeys := make([]*llrb.KeyInt, 0)
	fmt.Printf("number of elements {dict: %v, api:%v}\n", d.Len(), rb.Len())
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

//--------
// helpers
//--------

func compute_stats(cmd []interface{}, stats map[string]int) map[string]int {
	name := cmd[0].(string)
	count, ok := stats[name]
	if !ok {
		count = 0
	}
	stats[name] = count + 1
	return stats
}

func printStats(stats map[string]int) {
	keys, total := []string{}, 0
	for name, count := range stats {
		keys = append(keys, name)
		total += count
	}
	fmt.Printf("total commands : %v\n", total)
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%v: %v\n", key, stats[key])
	}
}

var writeOps = map[string]bool{
	"delmin": true,
	"delmax": true,
	"upsert": true,
	"insert": true,
	"delete": true,
}

func isReadOp(cmd []interface{}) bool {
	return !writeOps[cmd[0].(string)]
}
