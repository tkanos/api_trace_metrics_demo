package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var portConfigTests = []struct {
	key   string
	value string
	out   int
	env   string
	desc  string
}{
	{"PORT", "421", 421, "", "Customized configuration expect"},
	{"PORT", "", 5000, "DEV", "Default port configuration expect"},
}

//TestGet_Port Test the port configuration
func TestConfig(t *testing.T) {
	for _, test := range portConfigTests {
		// Arrange
		os.Setenv("ENVIRONMENT", test.env)
		os.Setenv(test.key, test.value)

		//Act
		Init()

		//Assert
		assert.Equal(t, test.out, Port, test.desc)
	}
}
