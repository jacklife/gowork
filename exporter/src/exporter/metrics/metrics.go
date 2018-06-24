package metrics

import (
	"sync"
	"github.com/prometheus/client_golang/prometheus"
	"database/sql"
	"github.com/prometheus/common/log"
)

const (
	pgNameSpace = "pg"
)

type PGMetrics struct {
	lock              sync.Mutex
	totalTableCount   prometheus.Gauge
	totalDbSize       prometheus.Gauge
	maxConnections    prometheus.Gauge
	curConnections    prometheus.Gauge
	activeConnections prometheus.Gauge
	dbSize            prometheus.GaugeVec
}

func NewPGMetrics(ip *string, port *string, version *string) (*PGMetrics) {
	var label = prometheus.Labels{"instance": *ip, "port": *port, "version": *version}
	return &PGMetrics{
		totalTableCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   pgNameSpace,
				Name:        "table_total_count",
				Help:        "Total tables count of current postgres instance",
				ConstLabels: label,
			},
		),
		totalDbSize: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   pgNameSpace,
				Name:        "total_database_size_bytes",
				Help:        "Total databases size of current postgres instance in bytes",
				ConstLabels: label,
			},
		),
		maxConnections: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   pgNameSpace,
				Name:        "connections_max_count",
				Help:        "Maximum connections of current postgres instance",
				ConstLabels: label,
			},
		),
		curConnections: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   pgNameSpace,
				Name:        "connections_current_count",
				Help:        "Connections established of current postgres instance",
				ConstLabels: label,
			},
		),
		activeConnections: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   pgNameSpace,
				Name:        "connections_active_count",
				Help:        "Busy established connections of current postgres instance",
				ConstLabels: label,
			},
		),
		dbSize: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   pgNameSpace,
				Name:        "database_size_bytes",
				Help:        "Size of a single database in bytes",
				ConstLabels: label,
			},
			[]string{"database"},
		),
	}
}


func (pgMetrics *PGMetrics) Scrape(db *sql.DB) float64 {

	pgMetrics.lock.Lock()
	defer pgMetrics.lock.Unlock()

	var errCount float64
	tableCount, err := queryTableCount(db)
	if err != nil {
		log.Error(err)
		errCount++
	}
	pgMetrics.totalTableCount.Set(tableCount)

	totalSize, err := queryTotalSize(db)
	if err != nil {
		log.Error(err)
		errCount++
	}
	pgMetrics.totalDbSize.Set(totalSize)

	maxConn, err := queryMaxConnections(db)
	if err != nil {
		log.Error(err)
		errCount++
	}
	pgMetrics.maxConnections.Set(maxConn)

	curConn, err := queryCurrentConnections(db)
	if err != nil {
		log.Error(err)
		errCount++
	}
	pgMetrics.curConnections.Set(curConn)

	activeConn, err := queryActiveConnections(db)
	if err != nil {
		log.Error(err)
		errCount++
	}
	pgMetrics.activeConnections.Set(activeConn)

	dbs, err := queryDBSize(db)
	if err != nil {
		log.Error(err)
		errCount++
	}
	for dbName, dbSize := range dbs {
		pgMetrics.dbSize.WithLabelValues(dbName).Set(dbSize)
	}

	return errCount

}

func (pgMetrics *PGMetrics) Describe(ch chan<- *prometheus.Desc) {

	pgMetrics.totalDbSize.Describe(ch)
	pgMetrics.totalTableCount.Describe(ch)
	pgMetrics.maxConnections.Describe(ch)
	pgMetrics.curConnections.Describe(ch)
	pgMetrics.activeConnections.Describe(ch)
	pgMetrics.dbSize.Describe(ch)

}

func (pgMetrics *PGMetrics) Collect(ch chan<- prometheus.Metric) {
	pgMetrics.totalTableCount.Collect(ch)

	pgMetrics.totalDbSize.Collect(ch)
	pgMetrics.maxConnections.Collect(ch)
	pgMetrics.curConnections.Collect(ch)
	pgMetrics.activeConnections.Collect(ch)
	pgMetrics.dbSize.Collect(ch)
}
