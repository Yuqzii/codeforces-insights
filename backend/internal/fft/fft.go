package fft

import (
	"math"
	"math/bits"
	"math/cmplx"
)

// @param a Slice to perform FFT on. Length must be a power of 2.
func FFT(a []complex128) []complex128 {
	return fftIterative(a)
}

// Iterative implementation of the Cooley-Tukey radix-2 algorithm.
// Based on the pseudocode here:
// https://en.wikipedia.org/wiki/Cooley%E2%80%93Tukey_FFT_algorithm#Data_reordering,_bit_reversal,_and_in-place_algorithms
// @param a The slice to perform FFT on. Length must be a power of 2.
func fftIterative(a []complex128) []complex128 {
	bitReverseCopy(a)

	for s := 1; s < bits.Len(uint(len(a))); s++ {
		m := 1 << s
		exponent := complex(0, -2*math.Pi/float64(m))
		omegaM := cmplx.Exp(exponent)

		for k := 0; k < len(a); k += m {
			omega := complex(1, 0)
			for j := 0; j < m/2; j++ {
				t := omega * a[k+j+m/2]
				u := a[k+j]
				a[k+j] = u + t
				a[k+j+m/2] = u - t
				omega *= omegaM
			}
		}
	}

	return a
}

// Works because IFFT(x) = 1/n * FFT(timeReverse(x)),
// where timeReverse is swapping elements symmetrically around the center excluding index 0.
// Exploits DFT symmetry to compute IFFT using forward FFT.
func IFFT(a []complex128) []complex128 {
	n := len(a)
	res := make([]complex128, n)
	copy(res, a)

	// Time reversal
	for i := 1; i < n/2; i++ {
		j := n - i
		res[i], res[j] = res[j], res[i]
	}

	res = FFT(res)

	// Scale output by 1/n
	invN := complex(1.0/float64(n), 0)
	for i := range res {
		res[i] *= invN
	}

	return res
}

// @param a Slice to do bit reversal on. Length must be power of 2. Modified in place.
func bitReverseCopy(a []complex128) {
	n := len(a)
	if n <= 2 {
		return
	}

	leading := bits.LeadingZeros(uint(n - 1))

	for i := range n {
		// After reversing all the leading zeros will be on the right side.
		// Therefore bitshift with that amount to keep it correctly inside the length of a.
		r := int(bits.Reverse(uint(i)) >> uint(leading))
		if i < r {
			a[i], a[r] = a[r], a[i]
		}
	}
}
