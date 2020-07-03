package logger

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-kit/kit/endpoint"
)

type transportDataLogging struct {
	Request      interface{} `json:"request"`
	Response     interface{} `json:"respons"`
	ResponseTime float64     `json:"response_time"`
}

//LoggingMiddleware transport logging middleware
func LoggingMiddleware(transportLogger *logrus.Logger) func(endpoint.Endpoint) endpoint.Endpoint {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			s := time.Now()
			resp, errResp := next(ctx, request)
			endtime := time.Now()
			var respLog interface{}
			if errResp != nil {
				respLog = errResp
			} else {
				respLog = resp
			}
			v := transportDataLogging{
				Response:     respLog,
				Request:      parseRequest(request),
				ResponseTime: endtime.Sub(s).Seconds(),
			}
			jsonLog, err := json.Marshal(v)
			if err != nil {
				log.Println(err.Error())
			}
			transportLogger.Debug(string(jsonLog))
			return resp, errResp
		}
	}
}

type httpRequestLog struct {
	RequestURI    string `json:"request_uri"`
	RequestMethod string `json:"request_method"`
	RemoteAddr    string `json:"remote_addr"`
}

func parseRequest(req interface{}) interface{} {
	switch v := req.(type) {
	case *http.Request:
		return httpRequestLog{
			RequestMethod: v.Method,
			RequestURI:    v.RequestURI,
			RemoteAddr:    v.RemoteAddr,
		}
	default:
		return v
	}
}
