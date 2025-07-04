package metrics

import (
	"flyhorizons-emailservice/internal/health"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMetricsRoutes(router *gin.Engine, rabbitMQCheck health.RabbitMQCheck) {
	rabbitMQHealthGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rabbitmq_health",
		Help: "RabbitMQ health status: 1 for up, 0 for down",
	})

	// Register both metrics
	prometheus.MustRegister(rabbitMQHealthGauge)

	go func() {
		for {
			if rabbitMQCheck.Pass() {
				rabbitMQHealthGauge.Set(1)
			} else {
				rabbitMQHealthGauge.Set(0)
			}

			time.Sleep(10 * time.Second) // adjust interval as needed
		}
	}()

	// Expose /metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
