package platform

import "github.com/prometheus/client_golang/prometheus"

type Prometheus struct {
	activeUserGauge       prometheus.Gauge
	uniqueActiveUserGauge prometheus.Gauge
}

func NewPrometheus() *Prometheus {
	activeUserGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_user",
			Help: "Number of active user",
		},
	)

	uniqueActiveUserGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "unique_active_user",
			Help: "Number of unique active user",
		},
	)

	prometheus.MustRegister(
		activeUserGauge,
		uniqueActiveUserGauge,
	)

	return &Prometheus{
		activeUserGauge:       activeUserGauge,
		uniqueActiveUserGauge: uniqueActiveUserGauge,
	}
}

func (p *Prometheus) GetActiveUserGauge() prometheus.Gauge {
	return p.activeUserGauge
}

func (p *Prometheus) GetUniqueActiveUserGauge() prometheus.Gauge {
	return p.uniqueActiveUserGauge
}
