package core

import (
	"testing"
)

func TestNewRange(t *testing.T) {
	r := NewRange(10)
	if r.Min != 10 && r.Max != 10 {
		t.Fatal("Max and min were not set correctly.")
	}

	r2 := NewRange(10)
	r2.Update(100)
	if r2.Min != 100 && r2.Max != 110 {
		t.Fatal("Failed to calculate range update correctly.")
	}

	r3 := NewRange(10)
	r3.Update(5)
	if r3.Min != 10 && r3.Max != 15 {
		t.Fatal("Failed to calculate range update correctly.")
	}
}
