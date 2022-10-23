package web

import (
	"context"
	"net/http"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/fruits"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPServer is a factory to create http servers for this project.
func NewHTTPServer(fruitEndpoints fruits.Endpoints, logger *loggers.Logger) http.Handler {
	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/fruit/{id}").Handler(
		httptransport.NewServer(
			fruitEndpoints.GetFruitWithIDEndpoint,
			makeDecodeGetFruitWithIDRequest(logger),
			makeEncodeGetFruitWithIDResponse(logger)),
	)
	router.Methods(http.MethodPut).Path("/fruit").Handler(
		httptransport.NewServer(
			fruitEndpoints.CreateFruitEndpoint,
			makeDecodeCreateFruitRequest(logger),
			makeEncodeCreateFruitRequest(logger)),
	)
	router.Methods(http.MethodGet).Path("/fruit").Handler(
		httptransport.NewServer(
			fruitEndpoints.SearchFruitsEndpoint,
			makeDecodeSearchFruitsRequest(logger),
			makeEncodeSearchFruitsResponse(logger)),
	)
	router.Methods(http.MethodGet).Path("/status").Handler(
		httptransport.NewServer(
			fruitEndpoints.GetStatusEndpoint,
			makeEmptyDecoder(logger),
			makeEncodeGetStatusResponse(logger)),
	)
	router.Methods(http.MethodGet).Path("/heartbeat").Handler(
		httptransport.NewServer(
			MakeGetHeartbeatEndpoint(logger),
			makeEmptyDecoder(logger),
			makeEncodeHeartbeatResponse(logger)),
	)

	return router
}

// MakeGetHeartbeatEndpoint service endpoint is a heartbeat.
func MakeGetHeartbeatEndpoint(logger *loggers.Logger) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		heartbeat := Result{
			Success: true,
		}

		logger.Debug(
			"get fruit heartbeat",
			loggers.Fields{
				"method": "GetHeartbeatEndpoint",
				"result": heartbeat,
			},
		)

		return heartbeat, nil
	}
}
