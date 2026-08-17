package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	ht "github.com/ogen-go/ogen/http"
	ogenjson "github.com/ogen-go/ogen/json"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"

	"cadence/api/server/publicapi"
	commonogenerrors "cadence/pkg/common/ogenerrors"
)

func ErrorHandler(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	var (
		ctError *validate.InvalidContentTypeError
		ogenErr ogenerrors.Error
		appErr  *commonogenerrors.Error
	)
	statusCode := http.StatusInternalServerError
	e := jx.GetEncoder()
	switch {
	case errors.Is(err, ht.ErrNotImplemented):
		statusCode = http.StatusNotImplemented
	case errors.As(err, &ctError):
		// Takes precedence over Error.
		statusCode = http.StatusUnsupportedMediaType
	case errors.As(err, &ogenErr):
		statusCode = ogenErr.Code()
		if statusCode == http.StatusBadRequest {
			resp := errorResponse(string(commonogenerrors.ErrCodeInvalidInput), ogenErr.Error(), nil)
			resp.Encode(e)
		}
	case errors.As(err, &appErr):
		statusCode = http.StatusBadRequest
		resp := errorResponse(string(appErr.Code), appErr.Message, appErr.Details)
		resp.Encode(e)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = e.WriteTo(w)
}

func errorResponse(code, message string, data map[string]interface{}) ogenjson.Marshaler {
	var optErrorData publicapi.OptErrorData
	if data != nil {
		errorData, _ := convertToRawMap(data)
		optErrorData = publicapi.NewOptErrorData(errorData)
	}
	return &publicapi.BadRequestResponseBody{
		Code:    publicapi.ErrorCode(code),
		Message: message,
		Data:    optErrorData,
	}
}

func convertToRawMap(input map[string]interface{}) (map[string]jx.Raw, error) {
	rawMap := make(map[string]jx.Raw)
	for key, val := range input {
		jsonBytes, err := json.Marshal(val)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal key %s", key)
		}
		rawMap[key] = jsonBytes
	}
	return rawMap, nil
}
