package commonmodel

import "strconv"

const FloatingPointPrecision = 1e9

type Concurrency int64

// Int64 returns the Concurrency value as an int64 (original, no scaling by FloatingPointPrecision).
func (c Concurrency) Int64() int64 {
	return int64(c)
}

func (c Concurrency) String() string {
	return strconv.FormatFloat(float64(c)/FloatingPointPrecision, 'f', -1, 64)
}

// Float64 returns the Concurrency value as a float64 but scaled by FloatingPointPrecision.
func (c Concurrency) Float64() float64 {
	return float64(c) / FloatingPointPrecision
}

// NewConcurrency creates a new Concurrency instance from a float64 value.
// The value is multiplied by FloatingPointPrecision to maintain precision.
func NewConcurrency(value float64) Concurrency {
	return Concurrency(int64(value * FloatingPointPrecision))
}
