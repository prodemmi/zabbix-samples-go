package main

import "github.com/prometheus/client_golang/prometheus"

type UserMetrics struct {
	CountOnline    prometheus.Gauge
	CountMaxOrder  prometheus.Gauge
	MaxPaymentByIP prometheus.GaugeVec
}

func (m UserMetrics) UpdateOnlineCount(count int64) {
	m.CountOnline.Set(float64(count))
}

func RegisterUserMetrics() *UserMetrics {
	metrics := &UserMetrics{
		CountOnline: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "user",
			Name:      "count_online",
			Help:      "Number of currently online users.",
		}),
		CountMaxOrder: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "user",
			Name:      "count_max_order",
			Help:      "Number of currently max order users.",
		}),
		MaxPaymentByIP: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "user",
			Name:      "max_payment_by_ip",
			Help:      "Get max payments by ip",
		}, []string{"country"}),
	}
	prometheus.MustRegister(metrics.CountOnline, metrics.CountMaxOrder, metrics.MaxPaymentByIP)

	return metrics
}
