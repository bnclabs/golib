package golib

import "testing"

func TestAverage(t *testing.T) {
	a := Average{}
	for i := 1; i <= 100; i++ {
		a.Add(float64(i))
	}
	if a.Count() != 100 {
		t.Fatalf("Count() mismatch!")
	}
	sum := ((100.0 * 101.0) / 2.0)
	if a.Sum() != sum {
		t.Fatalf("Sum() mismatch!")
	}
	if a.Mean() != sum/float64(a.Count()) {
		t.Fatalf("Mean() mismatch!")
	}
	if a.Variance() != 833.25 {
		t.Fatalf("Variance() mismatch!")
	}
	if a.Sd() != 28.86607004772212 {
		t.Fatalf("Sd() mismatch!")
	}
}
