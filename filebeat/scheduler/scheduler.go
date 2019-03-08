package scheduler

import (
	"github.com/elastic/beats/libbeat/logp"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

func Start() {
	gcInterval, ok := os.LookupEnv("GC_INTERVAL")
	_, debugOk := os.LookupEnv("GC_INTERVAL_DEBUG")
	if ok {
		interval, parseError := time.ParseDuration(gcInterval)
		if parseError != nil {
			logp.Err("Error parsing GC_INTERVAL variable [%s]", gcInterval)
		} else {
			logp.Info("Starting FreeOSMemory scheduler using interval [%s]", gcInterval)
			startTicker(interval, debugOk)
		}
	} else {
		logp.Info("Variable GC_INTERVAL is not set. FreeOSMemory scheduler not started.")
	}
}

func startTicker(interval time.Duration, debugOk bool) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			debug.FreeOSMemory()
			if debugOk {
				printMemUsage()
			}
		}
	}()
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	logp.Info("Memory stats: Alloc = %v MiB \tTotalAlloc = %v MiB \tSys = %v MiB \tNumGC = %v", bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
