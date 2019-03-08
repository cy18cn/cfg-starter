package handlers

import (
	"github.com/cy18cn/zlog"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"net/http"
)

func loggingHandler(next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		traceId := xid.New().String()
		method := request.Method
		zlog.Info(traceId,
			zap.String("uri", request.RequestURI),
			zap.String("method", request.Method))

		var contentType string
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			contentType, _ = getRequestContentType(request)
		}

		switch {
		case contentType == "application/json":
			zlog.Info(traceId,
				zap.String("contentType", contentType),
				zap.String("body", request.Form["body"][0]))
		default:
			zlog.Info(traceId,
				zap.String("contentType", contentType),
				zap.Any("params", request.Form))
		}
		request.Form["traceId"] = []string{traceId} // add traceId for logging
		next(writer, request, params)
	}
}
