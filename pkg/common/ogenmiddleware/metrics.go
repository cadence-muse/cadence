package ogenmiddleware

import (
	"time"

	"github.com/ogen-go/ogen/middleware"

	"cadence/pkg/common/metrics"
)

func NewMetricsMiddleware(m *metrics.Metrics) middleware.Middleware {
	return func(request middleware.Request, next middleware.Next) (middleware.Response, error) {
		start := time.Now()
		resp, err := next(request)
		duration := time.Since(start).Seconds()

		status := "success"
		if err != nil {
			status = "error"
		}

		m.HTTPRequestsTotal.WithLabelValues(request.OperationID, status).Inc()
		m.HTTPRequestDurationSec.WithLabelValues(request.OperationID).Observe(duration)

		return resp, err
	}
}
