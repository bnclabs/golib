package main

import "fmt"
import "log"
import "os"
import "strings"
import "flag"
import "time"
import "math/rand"
import "runtime/pprof"

import "github.com/prataprc/golib/llrb"
import "github.com/prataprc/golib"

var options struct {
	bcount  int
	readers int
	pprof   string
	mprof   string
	ops     []string
	algo    map[string]llrb.MemStore
	mtick   int
}

var insInts = make([]llrb.Item, 0)
var upsInts = make([]llrb.Item, 0)
var getInts = make([]llrb.Item, 0)
var delInts = make([]llrb.Item, 0)

func argParse() {
	var ops string
	var algo string
	flag.IntVar(&options.bcount, "count", 1,
		"data set size")
	flag.IntVar(&options.readers, "readers", 2,
		"number of concurrent readers to use with MVCC")
	flag.StringVar(&options.pprof, "pprof", "",
		"filename to save pprof o/p")
	flag.StringVar(&options.mprof, "mprof", "",
		"filename to save mprof o/p")
	flag.StringVar(&ops, "ops", "",
		"operations to profile")
	flag.StringVar(&algo, "algo", "llrb",
		"operations to profile")
	flag.IntVar(&options.mtick, "mtick", 0,
		"periodic tick to dump mem-stat, in mS")
	flag.Parse()

	options.ops = make([]string, 0)
	for _, op := range strings.Split(ops, ",") {
		if strings.Trim(op, " ") == "" {
			continue
		}
		options.ops = append(options.ops, op)
	}
	options.algo = make(map[string]llrb.MemStore)
	for _, algo := range strings.Split(algo, ",") {
		if strings.Trim(algo, " ") == "" {
			continue
		}
		switch algo {
		case "llrb":
			options.algo[algo] = llrb.NewLLRB()
		case "mvcc":
			options.algo[algo] = llrb.NewLLRBMVCC(10)
		}
	}
}

func main() {
	argParse()
	if options.mtick > 0 {
		go golib.MemProfiler(options.mtick, os.Stdout)
	}
	initialize()
	if len(options.ops) == 0 {
		return
	}

	startCPUProfile(options.pprof)
	count := options.bcount
	template := map[string][]interface{}{
		"min":    []interface{}{benchMin, count},
		"max":    []interface{}{benchMax, count},
		"get":    []interface{}{benchGet, len(getInts)},
		"upsert": []interface{}{benchUpsert, len(upsInts)},
		"range":  []interface{}{benchRange, count},
		"delete": []interface{}{benchDelete, len(delInts)},
	}
	for _, op := range options.ops {
		if op == "" {
			continue
		}
		fn := template[op][0].(func(llrb.MemStore))
		count := template[op][1].(int)
		for name, algo := range options.algo {
			fmsg := name + "." + op + " : %6d ns\n"
			timeit(fmsg, func() { fn(algo) }, count)
		}
	}
	takeMEMProfile(options.mprof)
	pprof.StopCPUProfile()
}

func benchMin(s llrb.MemStore) {
	for i := 0; i < options.bcount; i++ {
		s.Min()
	}
}
func benchMax(s llrb.MemStore) {
	for i := 0; i < options.bcount; i++ {
		s.Max()
	}
}
func benchInsert(s llrb.MemStore) {
	for _, key := range insInts {
		s.Insert(key)
	}
}
func benchUpsert(s llrb.MemStore) {
	for _, key := range upsInts {
		s.Upsert(key)
	}
}
func benchGet(s llrb.MemStore) {
	for _, key := range getInts {
		s.Get(key)
	}
}
func benchRange(s llrb.MemStore) {
	for i := 0; i < options.bcount; i++ {
		s.Range(nil, nil, "both", nil)
	}
}
func benchDelete(s llrb.MemStore) {
	for _, key := range delInts {
		s.Delete(key)
	}
}

func timeit(fmsg string, fn func(), count int) {
	now := time.Now()
	fn()
	fmt.Printf(fmsg, time.Since(now).Nanoseconds()/int64(count))
}

func initialize() {
	if len(options.ops) == 0 {
		startCPUProfile(options.pprof)
	}
	for name, store := range options.algo {
		for i := 0; i < options.bcount; i++ {
			item := &llrb.KeyInt{int64(i), time.Now().UnixNano()}
			store.Insert(item)
		}
		fmt.Printf("inserted %d items into %s\n", store.Len(), name)
	}

	if len(options.ops) > 0 {
		for i := 0; i < options.bcount; i++ {
			x, y := rand.Intn(options.bcount), time.Now().UnixNano()
			getInts = append(getInts, &llrb.KeyInt{int64(x), y})
			insInts = append(insInts, &llrb.KeyInt{int64(x + options.bcount), y})
			upsInts = append(upsInts, &llrb.KeyInt{int64(x), y})
			delInts = append(delInts, &llrb.KeyInt{int64(x * (rand.Intn(2) + 1)), y})
		}
		fmt.Printf("generated %v insert keys\n", len(insInts))
		fmt.Printf("generated %v upsert keys\n", len(upsInts))
		fmt.Printf("generated %v get keys\n", len(getInts))
		fmt.Printf("generated %v del keys\n", len(delInts))
	}
	if len(options.ops) == 0 {
		takeMEMProfile(options.mprof)
		pprof.StopCPUProfile()
	}
}

func startCPUProfile(filename string) {
	if filename == "" {
		return
	}
	fd, err := os.Create(filename)
	if err != nil {
		log.Fatalf("unable to create %q: %v\n", filename, err)
		return
	}
	pprof.StartCPUProfile(fd)
}

func takeMEMProfile(filename string) {
	if filename == "" {
		return
	}
	fd, err := os.Create(filename)
	if err != nil {
		log.Fatalf("unable to create %q: %v\n", filename, err)
	}
	pprof.WriteHeapProfile(fd)
	defer fd.Close()
}
