package atc2json

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalcChecksum(t *testing.T) {
	data := []byte{'A', 2, 3, 'z'}
	res := calcChecksum(data)
	assert.Equal(t, uint32(192), res, "192 and %d should be equal", res)
	assert.NotEqual(t, uint(0), res, "0 and %d should be not equal", res)
}

func TestCalcMillivolts(t *testing.T) {
	data := []int16{2000, 1000, 0, -1000, -2000}
	scale := float32(2000)
	res := calcMillivolts(data, scale)
	assert.Equal(t, []float32{1, 0.5, 0, -0.5, -1}, res, "Arrays should be equal")
}
