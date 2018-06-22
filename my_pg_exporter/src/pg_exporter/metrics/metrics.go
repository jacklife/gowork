package metrics

import (
	"sync"
	"github.com/prometheus/client_golang/prometheus"
)

type TableMetrics struct {
	mutex    sync.Mutex
	names    []string
	namesMap map[string]struct{}
	metrics  map[string]*prometheus.GaugeVec
}

type PGMetrics struct {
	mutex sync.Mutex
	totalCount prometheus.Gauge
	
}

