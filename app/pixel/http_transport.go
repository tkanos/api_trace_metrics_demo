package pixel

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// ErrInvalidBody thrown when the body of a request can not be parsed
var ErrInvalidBody = errors.New("invalid body")

// MakeHTTPHandler returns all http handler for the px service
func MakeHTTPHandler(endpoints Endpoints, tracer stdopentracing.Tracer) http.Handler {
	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	getpxHandler := kithttp.NewServer(
		endpoints.GetByID,
		decodeGetpxRequest,
		encodeResponse,
		append(options, kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "calling HTTP GET /{id}", nil)))...,
	)

	r := mux.NewRouter().PathPrefix("/pixel/").Subrouter().StrictSlash(true)

	r.Handle("/{id}", getpxHandler).Methods("GET")

	return r
}

func decodeGetpxRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)

	return GetpxRequest{ID: vars["id"]}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch err {
	case ErrInvalidBody:
		w.WriteHeader(http.StatusBadRequest)
	case ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		//logger.Error("", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
