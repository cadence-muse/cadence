package ogenmiddleware

import (
	"time"

	"github.com/ogen-go/ogen/middleware"

	"cadence/pkg/common/log"
)

func NewLoggingMiddleware(logger log.Logger) middleware.Middleware {
	return func(request middleware.Request, next middleware.Next) (middleware.Response, error) {
		start := time.Now()
		resp, err := next(request)
		duration := time.Since(start).String()

		fields := log.Fields{
			"args":     getParamsForLog(request),
			"duration": duration,
			"method":   request.OperationID,
		}

		loggerWithFields := logger.WithFields(fields)
		if err == nil {
			loggerWithFields.Info("call finished")
		} else {
			loggerWithFields.Error(err, "call failed")
		}

		return resp, err
	}
}

func getParamsForLog(request middleware.Request) any {
	var params map[string]any
	if len(request.Params) > 0 {
		params = make(map[string]any)
		for param, value := range request.Params {
			params[param.Name] = value
		}
	}
	options := getTrimForLogOptions()
	result := make(map[string]any)
	if len(params) > 0 {
		result["params"] = log.TrimForLogs(params, options)
	}
	if request.Body != nil {
		result["body"] = log.TrimForLogs(request.Body, options)
	}
	return result
}

func getTrimForLogOptions() log.TrimForLogOptions {
	trimForLogsOptions := log.DefaultTrimForLogsOpts
	trimForLogsOptions.SensitiveFields = []string{
		"Password",
		"password",
		"newPassword",
		"NewPassword",
		"secret",
		"authorization",
		"auth",
	}
	trimForLogsOptions.SensitivePlaceholder = "HIDDEN"
	return trimForLogsOptions
}
