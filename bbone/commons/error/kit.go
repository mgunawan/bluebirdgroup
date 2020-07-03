package error

import (
	"context"
	"encoding/json"
	"net/http"

	"bbone/commons/constant"
)

//DefaultJSONEncodeError ...
func DefaultJSONEncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set(constant.HTTPHeaderContentType, constant.ApplicationJSONUTF8Type)
	code := http.StatusInternalServerError
	message := "Something went wrong, please contact your administrator"
	errorCode := http.StatusText(code)
	if he, ok := err.(HTTPError); ok {
		code = he.Code
		message = he.Message
		errorCode = he.ErrorCode
	} else {
		message = err.Error()
	}
	if message != "write /dev/stdout: input/output error" {
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(HTTPError{
			Code:      code,
			Message:   message,
			ErrorCode: errorCode,
		})
	}
}
