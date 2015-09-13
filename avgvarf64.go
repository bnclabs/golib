package golib

import "math"

// AverageF64 maintains the average and variance of a stream
// of numbers in a space-efficient manner.
type AverageF64 struct {
	count int64
	sum   float64
	sumsq float64
}

// Add a sample to counting average.
func (av *AverageF64) Add(sample float64) {
	av.count++
	av.sum += sample
	av.sumsq += sample * sample
}

// GetCount return the number of samples counted so far.
func (av *AverageF64) Count() int64 {
	return av.count
}

// Mean return the sum of all samples by number of samples so far.
func (av *AverageF64) Mean() float64 {
	return av.sum / float64(av.count)
}

// GetTotal return the sum of all samples so far.
func (av *AverageF64) Sum() float64 {
	return av.sum
}

// Variance return the variance of all samples so far.
func (av *AverageF64) Variance() float64 {
	a := av.Mean()
	return av.sumsq/float64(av.count) - a*a
}

// GetStdDev return the standard-deviation of all samples so far.
func (av *AverageF64) Sd() float64 {
	return math.Sqrt(av.Variance())
}
