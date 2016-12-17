package gwu

import "testing"

func TestID(ts *testing.T) {
	for i := ID(0); i < 512; i++ {
		k, e := AtoID(i.String())
		if e != nil || k != i {
			ts.Fatal("ID error", i, k)
		}
	}
	for i := ^ID(0); i > ^ID(0)-512; i-- {
		k, e := AtoID(i.String())
		if e != nil || k != i {
			ts.Fatal("ID error", i, k)
		}
	}
}
