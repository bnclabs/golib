package llrb

import "math"

// Average maintains the average and variance of a stream
// of numbers in a space-efficient manner.
type Average struct {
	count      int64
	sum, sumsq float64
}

// Add a sample to counting average.
func (av *Average) Add(sample float64) {
	av.count++
	av.sum += sample
	av.sumsq += sample * sample
}

// GetCount return the number of samples counted so far.
func (av *Average) GetCount() int64 { return av.count }

// GetAvg return the sum of all samples by number of samples so far.
func (av *Average) GetAvg() float64 { return av.sum / float64(av.count) }

// GetTotal return the sum of all samples so far.
func (av *Average) GetTotal() float64 { return av.sum }

// GetVar return the variance of all samples so far.
func (av *Average) GetVar() float64 {
	a := av.GetAvg()
	return av.sumsq/float64(av.count) - a*a
}

// GetStdDev return the standard-deviation of all samples so far.
func (av *Average) GetStdDev() float64 { return math.Sqrt(av.GetVar()) }
