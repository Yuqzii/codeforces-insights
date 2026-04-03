package fft

import (
	"math"
	"math/bits"
)

const MAXN = 1 << 16

var (
	forwardTwiddles = computeTwiddles(MAXN, false)
	inverseTwiddles = computeTwiddles(MAXN, true)
)

// Iterative implementation of the Cooley-Tukey DIF algorithm.
// Based on the pseudocode here:
// https://en.wikipedia.org/wiki/Cooley%E2%80%93Tukey_FFT_algorithm#Data_reordering,_bit_reversal,_and_in-place_algorithms
// @param a The slice to perform FFT on. Length must be a power of 2.
// Transform is done in-place and modifies this slice.
func fftDIF(a []complex128) {
	for s := bits.Len(uint(len(a))) - 1; s >= 1; s-- {
		m := 1 << s
		stride := MAXN / m

		for k := 0; k < len(a); k += m {
			for j := 0; j < m/2; j++ {
				twiddle := forwardTwiddles[j*stride]

				u := a[k+j]
				v := a[k+j+m/2]

				a[k+j] = u + v
				a[k+j+m/2] = (u - v) * twiddle
			}
		}
	}
}

// Iterative implementation of the Cooley-Tukey DIT algorithm.
// @param a The slice to perform IFFT on. Length must be a power of 2.
// Transform is done in-place and modifies this slice.
func ifftDIT(a []complex128) {
	for s := 1; s < bits.Len(uint(len(a))); s++ {
		m := 1 << s
		stride := MAXN / m

		for k := 0; k < len(a); k += m {
			for j := 0; j < m/2; j++ {
				twiddle := inverseTwiddles[j*stride]

				t := twiddle * a[k+j+m/2]
				u := a[k+j]

				a[k+j] = u + t
				a[k+j+m/2] = u - t
			}
		}
	}

	// Scale output by 1/n to return to original magnitude.
	invN := complex(1.0/float64(len(a)), 0)
	for i := range a {
		a[i] *= invN
	}
}

func computeTwiddles(n int, inverse bool) []complex128 {
	twiddles := make([]complex128, n/2)
	sign := -1.0
	if inverse {
		sign = 1.0
	}

	for i := range n / 2 {
		theta := sign * 2.0 * math.Pi * float64(i) / float64(n)
		s, c := math.Sincos(theta)
		twiddles[i] = complex(c, s)
	}

	return twiddles
}
