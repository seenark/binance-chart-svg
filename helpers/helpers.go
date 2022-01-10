package helpers

import (
	"fmt"
	"runtime"
	"time"
)

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func TimeLast24H() (start, end int) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -1)
	start = ToMilliseconds(startDate)
	end = ToMilliseconds(now)
	return
}

func ToMilliseconds(t time.Time) int {
	return int(t.UnixNano()) / 1e6
}
