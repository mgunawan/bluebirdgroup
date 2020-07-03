package error

import (
	"encoding/json"
	"net/http"
)

//HTTPError logistic error
type HTTPError struct {
	Internal  error  `json:"internal,omitempty"`
	ErrorCode string `json:"errorCode,omitempty"`
	Code      int    `json:"code"`
	Message   string `json:"message,empty"`
}

//NewHTTPError create new http error
func NewHTTPError(code int, errorCode, message string) HTTPError {
	if message == "" {
		message = http.StatusText(code)
	}

	return HTTPError{
		Code:      code,
		ErrorCode: errorCode,
		Message:   message,
	}
}

func (h HTTPError) Error() string {
	strByte, _ := json.Marshal(h)
	return string(strByte)
}

//SetInternal set internal error
func (h HTTPError) SetInternal(err error) HTTPError {
	h.Internal = err
	return h
}
