//+build unit

package pixel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MakeGetByIDEndpoint(t *testing.T) {
	fakeService := new(mockedService)
	fakeService.On("Getpx", "1").Return(&(px{}), nil)

	e := MakeGetByIDEndpoint(fakeService)
	u, err := e(nil, GetpxRequest{ID: "1"})

	assert.NotNil(t, e)
	assert.Nil(t, err)
	assert.NotNil(t, u)
}

func Test_MakeCreateEndpoint(t *testing.T) {
	fakeService := new(mockedService)
	fakeService.On("Createpx", px{}).Return(*new(string), nil)

	e := MakeCreateEndpoint(fakeService)
	_, err := e(nil, CreatepxRequest{})

	assert.NotNil(t, e)
	assert.Nil(t, err)
}