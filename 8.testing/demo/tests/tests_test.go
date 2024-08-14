package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type dtype struct {
	a int
}

var t dtype

func TestSomething(t *testing.T) {
	assert := assert.New(t)

	// assert equality
	assert.Equal(123, 123, "they should be equal")

	// assert inequality
	assert.NotEqual(123, 456, "they should not be equal")

	// assert for nil (good for errors)
	// assert.Nil(t)

	// // assert for not nil (good when you expect something)
	if assert.NotNil(t) {

		// now we know that object isn't nil, we are safe to make
		// further assertions without causing any errors
		assert.NotEqual("Something", t)
	}
}

func TestStatusNotDown(t *testing.T) {
	assert.NotEqual(t, "status", "down")
}

func TestCalculate(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		input    int
		expected int
	}{
		{2, 4},
		{-1, 1},
		{0, 2},
		{-5, -3},
		{99999, 100001},
	}

	for _, test := range tests {
		assert.Equal(Calculate(test.input), test.expected)
	}
}

func Calculate(x int) (result int) {
	result = x + 2
	return result
}
