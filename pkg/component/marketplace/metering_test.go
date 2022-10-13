package marketplace

import (
	"testing"
)

func TestQuantityRoundHalfDown(t *testing.T) {
	mkt := awsMeteringService{}
	qnt, err := mkt.round(float32(3.4))
	if err != nil {
		t.Error(err)
	}
	if *qnt != 3 {
		t.FailNow()
	}
}

func TestQuantityRoundHalf(t *testing.T) {
	mkt := awsMeteringService{}
	qnt, err := mkt.round(float32(3.5))
	if err != nil {
		t.Error(err)
	}
	if *qnt != 4 {
		t.FailNow()
	}
}

func TestQuantityRoundHalfUp(t *testing.T) {
	mkt := awsMeteringService{}
	qnt, err := mkt.round(float32(3.6))
	if err != nil {
		t.Error(err)
	}
	if *qnt != 4 {
		t.FailNow()
	}
}
