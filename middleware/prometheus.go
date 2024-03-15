package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "shopifyx_request",
		Help:    "Histogram of the /shopifyx request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10), // Adjust bucket sizes as needed
	}, []string{"path", "method", "status"})
)

func NewRoute(c *echo.Echo, route *echo.Route, handler echo.HandlerFunc) {
	c.Add(route.Method, route.Path, wrapHandlerWithMetrics(route.Path, route.Method, handler))
}

func wrapHandlerWithMetrics(path string, method string, handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()

		// Execute the actual handler and catch any errors
		err := handler(c)

		// Regardless of whether an error occurred, record the metrics
		duration := time.Since(startTime).Seconds()
		statusCode := fmt.Sprintf("%d", c.Response().Status)

		if err != nil {
			if c.Response().Status == http.StatusOK || c.Response().Status == http.StatusCreated { // Default status code
				c.Response().Status = http.StatusInternalServerError // Assume internal server error if not set
			}
			c.String(http.StatusInternalServerError, err.Error()) // Ensure the response reflects the error
		}

		RequestHistogram.WithLabelValues(path, method, statusCode).Observe(duration)
		return err
	}
}
