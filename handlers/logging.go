package handlers

import (
	"github.com/cy18cn/zlog"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/xid"
	"net/http"
)

func loggingHandler(next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		traceId := xid.New().String()
		method := request.Method
		zlog.Infof("(trace %s) handle for request url:%s, method: %s",
			traceId,
			request.RequestURI)

		var contentType string
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			contentType, _ = getRequestContentType(request)
		}

		switch {
		case contentType == "application/json":
			zlog.Infof("(trace %s) request content-type:%s, params: %v",
				traceId,
				contentType,
				request.Form["body"])
		default:
			zlog.Infof("(trace %s) request content-type:%s, params:",
				traceId,
				contentType,
				request.Form)
		}
		request.Form["traceId"] = []string{traceId} // add traceId for logging
		next(writer, request, params)
	}
}
