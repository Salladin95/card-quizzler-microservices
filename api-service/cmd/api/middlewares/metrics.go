package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"strings"
	"time"
)

var RequestMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "api_service",
	Subsystem:  "http",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"url", "status"})

func observeRequest(d time.Duration, status int, url string) {
	RequestMetrics.WithLabelValues(url, strconv.Itoa(status)).Observe(d.Seconds())
}

func MetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		url := c.Request().URL.String()
		// skip requests to metrics endpoints
		if strings.Contains(url, "metrics") {
			return next(c)
		}

		start := time.Now()
		err := next(c)

		status := c.Response().Status

		observeRequest(time.Since(start), status, url)
		return err
	}
}
