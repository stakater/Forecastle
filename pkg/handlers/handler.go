package handlers

import (
	"net/http"

	"github.com/stakater/Forecastle/pkg/log"
)

var (
	logger = log.New()
)

type Handler interface {
	Handle(responseWriter http.ResponseWriter, request *http.Request)
}
