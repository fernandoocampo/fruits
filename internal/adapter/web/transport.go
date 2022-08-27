package web

import (
	"net/http"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/fruits"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPServer is a factory to create http servers for this project.
func NewHTTPServer(endyear fruits.Endyear, logger *loggers.Logger) http.Handler {
	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/fruit/{id}").Handler(
		httptransport.NewServer(
			endyear.GetFruitWithIDEndpoint,
			makeDecodeGetFruitWithIDRequest(logger),
			makeEncodeGetFruitWithIDResponse(logger)),
	)
	router.Methods(http.MethodPut).Path("/fruit").Handler(
		httptransport.NewServer(
			endyear.CreateFruitEndpoint,
			makeDecodeCreateFruitRequest(logger),
			makeEncodeCreateFruitRequest(logger)),
	)
	router.Methods(http.MethodGet).Path("/fruit").Handler(
		httptransport.NewServer(
			endyear.SearchFruitsEndpoint,
			makeDecodeSearchFruitsRequest(logger),
			makeEncodeSearchFruitsResponse(logger)),
	)
	router.Methods(http.MethodGet).Path("/status").Handler(
		httptransport.NewServer(
			endyear.GetStatusEndpoint,
			makeEmptyDecoder(logger),
			makeEncodeGetStatusResponse(logger)),
	)

	return router
}
