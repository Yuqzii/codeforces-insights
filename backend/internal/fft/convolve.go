package fft

import (
	"math/bits"
	"slices"
)

func Convolve(f, g []complex128) []complex128 {
	// Pad f and g to the next power of 2.
	n := nextPow2(max(len(f), len(g)))
	f = slices.Grow(f, n-len(f))
	f = f[:n]
	g = slices.Grow(g, n-len(g))
	g = g[:n]

	fftDIF(f)
	fftDIF(g)

	for i := range f {
		f[i] *= g[i]
	}

	ifftDIT(f)

	return f
}

// Returns smallest power of 2 that is >= n.
func nextPow2(n int) int {
	if n == 0 {
		return 1
	}
	return 1 << uint64(bits.Len64(uint64(n-1)))
}
