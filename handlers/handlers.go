package handlers

import (
	"github.com/cy18cn/zlog"
	ghandlers "github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"net/http"
)

func NewHandlers(next httprouter.Handle) httprouter.Handle {
	return errorHandler(loggingHandler(next))
}

func LoggingHandler(next httprouter.Handle) httprouter.Handle {
	return loggingHandler(next)
}

func ErrHandler(next httprouter.Handle) httprouter.Handle {
	return errorHandler(next)
}

func ParseFormHandler(logger *zap.Logger, next http.Handler) http.Handler {
	return &parseFormHandler{
		next: next,
		log:  logger,
	}
}

func HttpRouterMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return ParseFormHandler(zlog.GetLogger(), ghandlers.CompressHandler(next))
}
