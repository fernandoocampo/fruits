package web

import (
	"context"
	"log"
	"net/http"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/fruits"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

const homeContent = `<!DOCTYPE html>
<html>
   <head>
	  <title>fruits home</title>
   </head>
   <body style="background-color:green;">
	  <h1>Products</h1>
	  <p>fruits service.</p>
   </body>
</html>
`

// NewHTTPServer is a factory to create http servers for this project.
func NewHTTPServer(fruitEndpoints fruits.Endpoints, logger *loggers.Logger) http.Handler {
	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/home").Handler(home{})
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

type home struct{}

func (h home) ServeHTTP(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set("Content-Type", "text/html;charset=UTF-8")

	_, err := res.Write([]byte(homeContent))
	if err != nil {
		log.Println("unable to write home content", err)
	}
}
