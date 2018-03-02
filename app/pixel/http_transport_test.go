//+build unit

package pixel

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/tkanos/api_trace_metrics_demo/app/logger"

	"github.com/stretchr/testify/assert"
)

func Test_MakeHTTPHandler(t *testing.T) {
	h := MakeHTTPHandler(Endpoints{}, mocktracer.New())

	assert.NotNil(t, h)
}

func Test_DecodeGetpxRequest(t *testing.T) {
	expected := GetpxRequest{}
	r, _ := http.NewRequest("GET", "/pxs/1", nil)

	req, err := decodeGetpxRequest(context.Background(), r)

	assert.Nil(t, err)
	assert.Equal(t, expected, req)
}

func Test_EncodeResponse(t *testing.T) {
	response := struct {
		ID string
	}{"12345"}
	expected := "{\"ID\":\"12345\"}\n"

	w := httptest.NewRecorder()
	err := encodeResponse(context.Background(), w, response)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(w.Body)

	assert.Nil(t, err)
	assert.Equal(t, expected, string(body))
}

func Test_EncodeResponse_Should_Return_JSON_ContentType(t *testing.T) {
	expected := "application/json; charset=utf-8"

	w := httptest.NewRecorder()
	err := encodeResponse(context.Background(), w, nil)

	assert.Nil(t, err)
	assert.Equal(t, expected, w.Header().Get("Content-Type"))
}

func Test_EncodeError_Should_Return_JSON_ContentType(t *testing.T) {
	//Arrange
	logger.New()
	expected := "application/json; charset=utf-8"

	w := httptest.NewRecorder()

	//Act
	encodeError(context.Background(), errors.New("error"), w)

	//Assert
	assert.Equal(t, expected, w.Header().Get("Content-Type"))
}

func Test_EncodeError(t *testing.T) {
	err := errors.New("fake error")
	expected := "{\"error\":\"fake error\"}\n"

	w := httptest.NewRecorder()
	encodeError(context.Background(), err, w)
	body, err := ioutil.ReadAll(w.Body)

	assert.Nil(t, err)
	assert.Equal(t, expected, string(body))
}

func Test_EncodeError_Should_Correctly_Map_Error(t *testing.T) {
	var flagtests = []struct {
		in  error
		out int
	}{
		{errors.New("not handled error"), http.StatusInternalServerError},
		{ErrInvalidBody, http.StatusBadRequest},
		{ErrNotFound, http.StatusNotFound},
	}

	for _, tt := range flagtests {
		w := httptest.NewRecorder()
		encodeError(context.Background(), tt.in, w)

		assert.Equal(t, tt.out, w.Code)
	}
}
