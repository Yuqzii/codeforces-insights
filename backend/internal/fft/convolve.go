package fft

import (
	"math/bits"
	"slices"
)

func Convolve(f, g []complex128) []complex128 {
	// Pad f and g to next power of 2
	n := nextPow2(len(f) + len(g))
	f = slices.Grow(f, n-len(f))
	f = f[:cap(f)]
	g = slices.Grow(g, n-len(g))
	g = g[:cap(g)]

	x := FFT(f)
	y := FFT(g)

	for i := range x {
		x[i] *= y[i]
	}

	return IFFT(x)
}

// Returns smallest power of 2 that is >= n.
func nextPow2(n int) int {
	if n == 0 {
		return 1
	}
	return 1 << uint64(bits.Len64(uint64(n-1)))
}
