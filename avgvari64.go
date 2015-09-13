package golib

import "math"

// AverageI64 maintains the average and variance of a stream
// of numbers in a space-efficient manner.
type AverageI64 struct {
	count int64
	sum   float64
	sumsq float64
}

// Add a sample to counting average.
func (av *AverageI64) Add(sample int64) {
	av.count++
	av.sum += float64(sample)
	av.sumsq += float64(sample) * float64(sample)
}

// GetCount return the number of samples counted so far.
func (av *AverageI64) Count() int64 {
	return av.count
}

// Mean return the sum of all samples by number of samples so far.
func (av *AverageI64) Mean() int64 {
	return int64(av.sum / float64(av.count))
}

// GetTotal return the sum of all samples so far.
func (av *AverageI64) Sum() float64 {
	return av.sum
}

// Variance return the variance of all samples so far.
func (av *AverageI64) Variance() float64 {
	a := av.Mean()
	return (av.sumsq / float64(av.count)) - (float64(a) * float64(a))
}

// GetStdDev return the standard-deviation of all samples so far.
func (av *AverageI64) Sd() float64 {
	return math.Sqrt(av.Variance())
}
