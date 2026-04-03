package fft

import (
	"math"
	"math/bits"
	"math/cmplx"
)

// Iterative implementation of the Cooley-Tukey DIF algorithm.
// Based on the pseudocode here:
// https://en.wikipedia.org/wiki/Cooley%E2%80%93Tukey_FFT_algorithm#Data_reordering,_bit_reversal,_and_in-place_algorithms
// @param a The slice to perform FFT on. Length must be a power of 2.
// Transform is done in-place and modifies this slice.
func fftDIF(a []complex128) {
	for s := bits.Len(uint(len(a))) - 1; s >= 1; s-- {
		m := 1 << s
		exp := cmplx.Exp(complex(0, -2*math.Pi/float64(m)))

		for k := 0; k < len(a); k += m {
			twiddle := complex(1, 0)
			for j := 0; j < m/2; j++ {
				u := a[k+j]
				v := a[k+j+m/2]

				a[k+j] = u + v
				a[k+j+m/2] = (u - v) * twiddle

				twiddle *= exp
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
		exp := cmplx.Exp(complex(0, 2*math.Pi/float64(m))) // Positive for inverse.

		for k := 0; k < len(a); k += m {
			twiddle := complex(1, 0)
			for j := 0; j < m/2; j++ {
				t := twiddle * a[k+j+m/2]
				u := a[k+j]

				a[k+j] = u + t
				a[k+j+m/2] = u - t
				twiddle *= exp
			}
		}
	}

	// Scale output by 1/n to return to original magnitude.
	invN := complex(1.0/float64(len(a)), 0)
	for i := range a {
		a[i] *= invN
	}
}
