package handlers

import "net/http"

type Handler interface {
	Handle(responseWriter http.ResponseWriter, request *http.Request)
}
