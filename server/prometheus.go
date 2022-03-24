package server

import (
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/aclrys/prebid-server/config"
	metricsconfig "github.com/aclrys/prebid-server/metrics/config"
)

func newPrometheusServer(cfg *config.Configuration, metrics *metricsconfig.DetailedMetricsEngine) *http.Server {
	proMetrics := metrics.PrometheusMetrics

	if proMetrics == nil {
		glog.Fatal("Prometheus metrics configured, but a Prometheus metrics engine was not found. Cannot set up a Prometheus listener.")
	}
	return &http.Server{
		Addr: cfg.Host + ":" + strconv.Itoa(cfg.Metrics.Prometheus.Port),
		Handler: promhttp.HandlerFor(proMetrics.Gatherer, promhttp.HandlerOpts{
			ErrorLog:            loggerForPrometheus{},
			MaxRequestsInFlight: 5,
			Timeout:             cfg.Metrics.Prometheus.Timeout(),
		}),
	}
}

type loggerForPrometheus struct{}

func (loggerForPrometheus) Println(v ...interface{}) {
	glog.Warningln(v...)
}
