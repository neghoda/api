package models

import (
	"time"
)

// Health - describes service health status.
type Health struct {
	GoroutinesNum      int                  `json:"goroutines_num"`
	Memory             Memory               `json:"memory"`
	IPs                []string             `json:"ips"`
	ExternalConnection []ExternalConnection `json:"external_connection"`
}

type Memory struct {
	// AllocBytes is cumulative bytes allocated for heap objects.
	// increases as heap objects are allocated
	AllocBytes string `json:"alloc_bytes"`

	// SysBytes is the total bytes of memory obtained from the OS.
	//
	// SysBytes is the sum of the XSys fields below. Sys measures the
	// virtual address space reserved by the Go runtime for the
	// heap, stacks, and other internal data structures. It's
	// likely that not all of the virtual address space is backed
	// by physical memory at any given moment, though in general
	// it all was at some point.
	SysBytes string `json:"sys_bytes"`

	// AllHeapObjects is the cumulative count of heap objects allocated.
	AllHeapObjects uint64 `json:"all_heap_objects"`

	// LiveHeapObjects is the number of live objects: is AllHeapObjects - memstat.Frees.
	LiveHeapObjects uint64 `json:"live_heap_objects"`

	// NumGC is the number of garbage collections.
	NumGC int64 `json:"num_gc"`

	// LastGC is the time of last garbage collection
	LastGC time.Time `json:"last_gc"`

	// GCPauseTotal is the duration of all pauses for all collections
	GCPauseTotal time.Duration `json:"gc_pause_total"`
}

type ConnectionStatus string

type ExternalConnection struct {
	Name   string           `json:"name"`
	Status ConnectionStatus `json:"status"`
}
