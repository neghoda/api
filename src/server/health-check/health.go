package healthcheck

import (
	"github.com/erp/api/src/models"
	"github.com/erp/api/src/server/handlers"
	"net"
	"net/http"
	"runtime"
	"runtime/debug"
	"sync"
)

var (
	memStats     = &runtime.MemStats{}
	gcStats      = &debug.GCStats{}
	mu           = &sync.Mutex{}
	healthStatus = models.Health{}
)

// Health send actual service health status.
func Health(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	runtime.ReadMemStats(memStats)
	debug.ReadGCStats(gcStats)

	ips, err := getIPs()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, nil)
		return
	}

	healthStatus.IPs = ips

	healthStatus.GoroutinesNum = runtime.NumGoroutine()

	healthStatus.Memory = models.Memory{
		AllocBytes:      BytesToString(memStats.Alloc),
		SysBytes:        BytesToString(memStats.Sys),
		AllHeapObjects:  memStats.Mallocs,
		LiveHeapObjects: memStats.Mallocs - memStats.Frees,
		NumGC:           gcStats.NumGC,
		LastGC:          gcStats.LastGC,
	}

	healthStatus.ExternalConnection = []models.ExternalConnection{}

	if _, err := handlers.HandleRequest(r); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, nil)
		return
	}

	handlers.SendResponse(w, http.StatusOK, healthStatus)
}

func getIPs() (ips []string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	ips = make([]string, 0, len(addrs))
	for _, ip := range addrs {
		ips = append(ips, ip.String())
	}

	return
}
