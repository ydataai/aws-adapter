package marketplace

import (
	"testing"
)

func TestQuantityConversion(t *testing.T) {
	mkt := awsMeteringService{}
	qnt, err := mkt.convert(float32(3.124124))
	if err != nil {
		t.Error(err)
	}
	if *qnt != 3124124 {
		t.FailNow()
	}
	res := float32(3.124124) * 1000000
	if res != 3124124 {
		t.FailNow()
	}
}
