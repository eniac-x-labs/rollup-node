package metrics

import (
	"github.com/eniac-x-labs/rollup-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

const Namespace = "eip-4844-rollup"

type Metricer interface {
	TxMetricer
}

type Metrics struct {
	ns       string
	Registry *prometheus.Registry
	Factory  metrics.Factory
	TxMetrics
}

var _ Metricer = (*Metrics)(nil)

func NewMetrics(procName string) *Metrics {
	if procName == "" {
		procName = "default"
	}
	ns := Namespace + "_" + procName

	registry := NewRegistry()
	factory := metrics.With(registry)

	return &Metrics{
		ns:       ns,
		Registry: registry,
		Factory:  factory,

		TxMetrics: MakeTxMetrics(ns, factory),
	}
}

type RegistryMetricer interface {
	Registry() *prometheus.Registry
}

func NewRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())
	return registry
}
