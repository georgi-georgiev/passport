package passport

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	days := d / (time.Hour * 24)
	d -= days * time.Hour * 24
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%d days, %02d:%02d:%02dh", days, h, m, s)
}

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		panic(err)
	}
}

// healthcheck handler with custom callback for printing service specific additional info
func HandleHealthCheckWithCallback(healthCallback func() string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write([]byte(healthCallback()))
		if err != nil {
			panic(err)
		}
	}
}

func HandleRuntimeStats(w http.ResponseWriter, r *http.Request) {
	// GET requests are served currently
	if r.Method != "GET" {
		http.Error(w, "Not Implemented.", http.StatusNotImplemented)
		return
	}

	s := "Running since (server time): " + startTime.Local().Format("2006-01-02 15:04:05") + "\n" +
		"Running for: " + fmtDuration(time.Since(startTime))
	s += "\n\n"
	s += fmt.Sprintln("Active goroutines:", runtime.NumGoroutine())
	s += fmt.Sprintln("Logical CPUs available to the runtime (GOMAXPROCS):", runtime.GOMAXPROCS(0))

	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	s += "\n"
	s += fmt.Sprintf("Heap in use: %sMB\n", formatMBs(m.Alloc))
	s += fmt.Sprintf("Heap allocated but not in use: %sMB\n", formatMBs(m.HeapInuse-m.HeapAlloc+m.HeapIdle))
	s += fmt.Sprintln("Heap objects: ", m.HeapObjects)
	s += fmt.Sprintf("Stack: %sMB\n", formatMBs(m.StackInuse))
	s += fmt.Sprintf("Total runtime reserved virtual space: %sMB\n", formatMBs(m.Sys))

	_, err := w.Write([]byte(s))
	if err != nil {
		panic(err)
	}
}

func formatMBs(n uint64) string {
	mb := n / 1_000_000
	remainder := n % 1_000_000
	// round MBs to two digits after the dot
	remainder += 5_000
	return strconv.FormatUint(mb, 10) + "." + fmt.Sprintf("%06d", remainder)[:2]
}
