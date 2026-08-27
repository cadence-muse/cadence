package ogenmiddleware

import (
	"encoding/json"
	"time"

	"github.com/go-faster/jx"
	"github.com/nightnoryu/go-kita/log"
	ogenjson "github.com/ogen-go/ogen/json"
	"github.com/ogen-go/ogen/middleware"
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
			params[param.Name] = unwrapParamValue(value)
		}
	}
	options := getTrimForLogsOptions()
	result := make(map[string]any)
	if len(params) > 0 {
		result["params"] = log.TrimForLogs(params, options)
	}
	if request.Body != nil {
		result["body"] = log.TrimForLogs(request.Body, options)
	}
	return result
}

// unwrapParamValue converts ogen-generated optional param types (e.g. OptUUID) into plain
// JSON-safe values. Their MarshalJSON writes through jx and returns empty bytes when unset,
// which stdlib encoding/json (used by zap and log.TrimForLogs) rejects as invalid JSON.
func unwrapParamValue(value any) any {
	marshaler, ok := value.(ogenjson.Marshaler)
	if !ok {
		return value
	}

	e := jx.Encoder{}
	marshaler.Encode(&e)
	raw := e.Bytes()
	if len(raw) == 0 {
		return nil
	}

	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return value
	}
	return decoded
}

func getTrimForLogsOptions() log.TrimForLogsOptions {
	trimForLogsOptions := log.DefaultTrimForLogsOpts
	trimForLogsOptions.SensitiveFields = []string{
		"Password",
		"password",
		"currentPassword",
		"CurrentPassword",
		"newPassword",
		"NewPassword",
		"secret",
		"authorization",
		"auth",
	}
	trimForLogsOptions.SensitivePlaceholder = "HIDDEN"
	return trimForLogsOptions
}
