package metrics

import (
	"net"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartServer(r *prometheus.Registry, hostname string, port int) (string, error) {
	addr := net.JoinHostPort(hostname, strconv.Itoa(port))
	h := promhttp.InstrumentMetricHandler(
		r, promhttp.HandlerFor(r, promhttp.HandlerOpts{}),
	)
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(addr, h); err != nil {
		return "", err
	}
	return addr, nil
}
