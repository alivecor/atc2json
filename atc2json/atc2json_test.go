package atc2json

import (
	"testing"
)

func TestCalcChecksum(t *testing.T) {
	data := []byte{'A', 2, 3, 'z'}
	res := calcChecksum(data)
	if uint32(192) != res {
		t.Errorf("192 and %d should be equal", res)
	}
	//assert.Equal(t, uint32(192), res, "192 and %d should be equal", res)
	if uint32(0) != res {
		t.Errorf("0 and %d should be equal", res)
	}
	//assert.NotEqual(t, uint(0), res, "0 and %d should be not equal", res)
}

func TestCalcMillivolts(t *testing.T) {
	data := []int16{2000, 1000, 0, -1000, -2000}
	scale := float32(2000)
	res := calcMillivolts(data, scale)
	if len(res) != 5 {
		t.Errorf("len(res) should be 5")
	}
	expected := []float32{1, 0.5, 0, -0.5, -1}
	for i := 0; i < 5; i++ {
		if expected[i] != res[i] {
			t.Errorf("Arrays should be equal. res[%d] is %f, expected %f", i, res[i], expected[i])
		}
	}
}
