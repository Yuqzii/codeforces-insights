package fft

import (
	"math"
	"math/cmplx"
)

func FFT(x []complex128) []complex128 {
	return fftRecursive(x, len(x), 1)
}

// Simple implementation of the Cooley-Tukey radix-2 algorithm.
// Recursively splits the DFT into two smaller smaller DFTs. O(nlogn) time complexity.
// Possible optimization is removing explicit recursion.
// @param a Slice of complex numbers to transform. Length must be power of 2.
// @param n Current length of DFT to process. Must be power of 2.
// @param s Current step. First compute even indices and the odd indices, then combines these.
// This idea is extended recursively to give it the logarithmic time complexity.
func fftRecursive(x []complex128, n, s int) []complex128 {
	if n == 1 {
		return []complex128{x[0]}
	}

	// Even and odd indices (only actually even and odd in first recursion)
	even := fftRecursive(x, n/2, 2*s)
	odd := fftRecursive(x[s:], n/2, 2*s)

	exp := cmplx.Rect(1, -2*math.Pi/float64(n)) // Roots of unity (please don't ask me to explain).
	twiddle := complex(1, 0)                    // Twiddle factor, accumulates rotation.
	res := make([]complex128, n)

	for k := 0; k < n/2; k++ {
		t := twiddle * odd[k]
		res[k] = even[k] + t     // Set first half
		res[k+n/2] = even[k] - t // Set second half

		twiddle *= exp // twiddle becomes exp^1, exp^2, exp^3, etc.
	}
	return res
}

// Works because IFFT(x) = 1/n * FFT(timeReverse(x)),
// where timeReverse is swapping elements symmetrically around the center excluding index 0.
// Exploits DFT symmetry to compute IFFT using forward FFT.
func IFFT(x []complex128) []complex128 {
	n := len(x)
	res := make([]complex128, n)
	copy(res, x)

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
