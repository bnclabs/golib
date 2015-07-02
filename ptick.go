package golib

import "os"
import "runtime"
import "fmt"
import "time"
import "encoding/json"

func MemProfiler(period int, fd *os.File) {
	var mstat runtime.MemStats
	tick := time.Tick(time.Duration(period) * time.Millisecond)
	for {
		<-tick
		runtime.ReadMemStats(&mstat)
		m := map[string]interface{}{
			"alloc":       mstat.Alloc,
			"total_alloc": mstat.TotalAlloc,
			// Main allocation heap statistics.
			"numgc":       mstat.NumGC,
			"pause_total": mstat.PauseTotalNs,
			"avgns":       int64(float64(mstat.PauseTotalNs) / float64(mstat.NumGC)),
		}
		data, _ := json.Marshal(m)
		fmt.Printf("memstat: %s\n", string(data))
	}
}
