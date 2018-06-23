package main

import (
	"database/sql"
	"net/http"
	"sync"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"myexporter/metrics"

	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/prometheus/common/log"
	"strings"
	"time"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	PGExporterNS = "PGExporter"
)

var db *sql.DB

var (
	listenAddress         = kingpin.Flag("web.listen-address", "Address to listen on for web interface.").Default("0.0.0.0:9292").Envar("PG_EXPORTER_WEB_LISTEN_ADDRESS").String()
	metricPath            = kingpin.Flag("web.exposition-path", "Path under which to expose metrics.").Default("/metrics").Envar("PG_EXPORTER_WEB_TELEMETRY_PATH").String()
	disableDefaultMetrics = kingpin.Flag("web.disable-default-metrics", "Do not include default metrics.").Default("false").Envar("PG_EXPORTER_DISABLE_DEFAULT_METRICS").Bool()
	onlyDumpMaps          = kingpin.Flag("dumpmaps", "Do not run, simply dump the maps.").Bool()
	//pgUserName            = kingpin.Flag("pg.username", "user name for postgresQL to be monitoring").Envar("POSTGRES_USER_NAME").Required().String()
	//pgPassword            = kingpin.Flag("pg.password", "password of pgUserName").Envar("POSTGRES_PASSWORD").Required().String()
	//dbName                = kingpin.Flag("pg.dbName", "postgres dbName").Required().Envar("POSTGRES_DB_NAME").String()
	pgAddress       = kingpin.Flag("pg.server-address", "postgres server address").Default("localhost:5432").Envar("POSTGRES_SERVER_ADDRESS").String()
	passwordEncrypt = kingpin.Flag("pg.passwdEncrypt", "Is pgPassword encrypted").Default("true").Bool()
)

type PGExporter struct {
	lock             sync.Mutex
	pgMetrics        metrics.PGMetrics
	totalScrapes     prometheus.Counter
	duration, errors prometheus.Gauge
}

func NewPGExporter(ip *string, port *string, version *string) *PGExporter {
	return &PGExporter{
		pgMetrics: *metrics.NewPGMetrics(ip, port, version),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: PGExporterNS,
			Name:      "exporter_scrapes_total",
			Help:      "Current total postgres scrapes.",
		}),
		duration: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: PGExporterNS,
			Name:      "exporter_last_scrape_duration_seconds",
			Help:      "The last scrape duration.",
		}),
		errors: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: PGExporterNS,
			Name:      "exporter_last_scrape_error",
			Help:      "The last scrape error status.",
		}),
	}
}

func (pgExporter *PGExporter) Describe(ch chan<- *prometheus.Desc) {
	pgExporter.pgMetrics.Describe(ch)

	pgExporter.totalScrapes.Describe(ch)
	pgExporter.duration.Describe(ch)
	pgExporter.errors.Describe(ch)
}

func (pgExporter *PGExporter) Collect(ch chan<- prometheus.Metric) {
	pgExporter.lock.Lock()
	defer pgExporter.lock.Unlock()
	pgExporter.scrape()

	pgExporter.pgMetrics.Collect(ch)

	pgExporter.totalScrapes.Collect(ch)

	pgExporter.duration.Collect(ch)
	pgExporter.errors.Collect(ch)

}

func (pgExporter *PGExporter) scrape() {
	now := time.Now().UnixNano()

	errCont:=pgExporter.pgMetrics.Scrape(db)

	pgExporter.totalScrapes.Inc()
	pgExporter.errors.Add(errCont)
	pgExporter.duration.Set(float64(time.Now().UnixNano()-now) / 1000000000)
}

func main() {
	kingpin.Parse()

	var err error
	connStr := "user=postgres password=123456 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("error opening connection to database: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("error opening connection to database: ", err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)

	ip, port, version := getAddress(db)

	exporter := NewPGExporter(&ip, &port, &version)

	prometheus.MustRegister(exporter)
	

	http.Handle(*metricPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
<head><title>PostgreSQL exporter</title></head>
<body>
<h1>PostgreSQL exporter</h1>
<p><a href='` + *metricPath + `'>Metrics</a></p>
</body>
</html>
`))
	})

	log.Infof("Starting Server: %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func getAddress(db *sql.DB) (ip string, port string, version string) {
	kingpin.Parse()
	ip = strings.Split(*pgAddress, ":")[0]
	port = strings.Split(*pgAddress, ":")[1]
	version = metrics.QueryVersion(db)
	return ip, port, version
}
