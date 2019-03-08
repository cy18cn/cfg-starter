package handlers

import (
	"fmt"
	"github.com/cy18cn/micro-svc-common/util"
	"github.com/cy18cn/zlog"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"net/http"
)

type responseLogger struct {
	w        http.ResponseWriter
	status   int
	size     int
	respBody []byte
}

func (self *responseLogger) Header() http.Header {
	return self.w.Header()
}

func (self *responseLogger) Write(b []byte) (int, error) {
	size, err := self.w.Write(b)
	self.size += size
	self.respBody = util.BytesCombine(self.respBody, b)
	return size, err
}

func (self *responseLogger) WriteHeader(statusCode int) {
	self.status = statusCode
	self.w.WriteHeader(statusCode)
}

func (self *responseLogger) Status() int {
	return self.status
}

func (self *responseLogger) Size() int {
	return self.size
}

func (self *responseLogger) RespBody() []byte {
	return self.respBody
}

func (self *responseLogger) Flush() {
	f, ok := self.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}

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
		w := &responseLogger{
			w: writer,
		}
		next(w, request, params)
		zlog.Info(fmt.Sprintf("%s response", traceId),
			zap.Int("status", w.status),
			zap.Int("responseSize", w.size),
			zap.String("response", string(w.respBody)))
	}
}
