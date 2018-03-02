//+build unit

package pixel

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewpxTracing_Should_Create_New_Service_Instance(t *testing.T) {
	//Arrange
	fake := new(mockedService)

	//Act
	tracer := NewpxTracing(fake)

	//Assert
	assert.NotNil(t, tracer)
}

func Test_Getpx_Should_Go_Throught_The_Method(t *testing.T) {
	//Arrange
	fake := new(mockedService)
	iD := "123"
	fake.On("Getpx", iD).Return(&px{}, nil)
	tracer := NewpxTracing(fake)

	//Act
	c, err := tracer.Getpx(context.Background(), iD)

	//Assert
	assert.NotNil(t, c)
	assert.NoError(t, err)
}

func Test_Getpx_With_Error_Should_Go_Throught_The_Method(t *testing.T) {
	//Arrange
	fake := new(mockedService)
	iD := "123"
	errorExpected := errors.New("test")
	fake.On("Getpx", iD).Return(&px{}, errors.New("test"))
	tracer := NewpxTracing(fake)

	//Act
	c, err := tracer.Getpx(context.Background(), iD)

	//Assert
	assert.NotNil(t, c)
	assert.Equal(t, errorExpected.Error(), err.Error(), "Middleware pxs tracing")
}
