package handlers

import "github.com/julienschmidt/httprouter"

func NewHandlers(next httprouter.Handle) httprouter.Handle {
	return errorHandler(loggingHandler(next))
}

func LoggingHandler(next httprouter.Handle) httprouter.Handle {
	return loggingHandler(next)
}

func ErrHandler(next httprouter.Handle) httprouter.Handle {
	return errorHandler(next)
}
