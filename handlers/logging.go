package handlers

import (
	"github.com/cy18cn/zlog"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/xid"
	"mime"
	"net/http"
	"net/url"
)

func loggingHandler(next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		traceId := xid.New().String()
		zlog.Infof("(trace %s) handle for request url:%s",
			traceId,
			request.RequestURI)
		contentType := request.Header.Get("Content-Type")
		// RFC 7231, section 3.1.1.5 - empty type
		//   MAY be treated as application/octet-stream
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		contentType, _, _ = mime.ParseMediaType(contentType)
		switch {
		case contentType == "application/json":
			zlog.Infof("(trace %s) request content-type:%s, params: %v",
				traceId,
				contentType,
				request.Form["body"])
		case contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data":
			zlog.Infof("(trace %s) request content-type:%s, params:",
				traceId,
				contentType,
				request.Form)
		}
		request.Form["traceId"] = []string{traceId}	// add traceId for logging
		next(writer, request, params)
	}
}

func filterExtFormData(form url.Values) url.Values {
	res := make(url.Values)
	for k, v := range form {
		if k == "traceId" {
			continue
		}

		res[k] = v
	}
	return res
}