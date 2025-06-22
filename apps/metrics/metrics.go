package metrics

import (
	"fmt"
	"net/http"
	"videoanalytics/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Counter for total clicks per ad
var AdClicks = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "ad_total_clickCount",
		Help: "Total number of ad clicks",
	},
	[]string{"ad_id"},
)

// Counter for total errors (by endpoint)
var RequestErrors = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_errors_total",
		Help: "Total errors by endpoint",
	},
	[]string{"endpoint"},
)

// Request count per endpoint
var HttpRequests = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "API_requests_total",
		Help: "Number of requests by endpoint and method",
	},
	[]string{"endpoint", "method"},
)

// Expose metrics on :9090
func StartMetricsServer() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		lPort := config.GetConfig().Metrics.Port
		http.ListenAndServe(fmt.Sprintf(":%d", lPort), nil)
	}()
}
