package fft

import (
	"math"
	"math/bits"
	"sync"
)

type twiddleTable struct {
	forward []complex128
	inverse []complex128
}

var (
	cacheMu sync.RWMutex
	cache   = make(map[int]twiddleTable)
)

// Iterative implementation of the Cooley-Tukey DIF algorithm.
// Based on the pseudocode here:
// https://en.wikipedia.org/wiki/Cooley%E2%80%93Tukey_FFT_algorithm#Data_reordering,_bit_reversal,_and_in-place_algorithms
// @param a The slice to perform FFT on. Length must be a power of 2.
// Transform is done in-place and modifies this slice.
func fftDIF(a []complex128) {
	twiddles := getTwiddles(len(a))

	for s := bits.Len(uint(len(a))) - 1; s >= 1; s-- {
		m := 1 << s
		stride := len(a) / m
		half := m / 2

		for k := 0; k < len(a); k += m {
			for j := range half {
				twiddle := twiddles.forward[j*stride]

				u := a[k+j]
				v := a[k+j+half]

				a[k+j] = u + v
				a[k+j+half] = (u - v) * twiddle
			}
		}
	}
}

// Iterative implementation of the Cooley-Tukey DIT algorithm.
// @param a The slice to perform IFFT on. Length must be a power of 2.
// Transform is done in-place and modifies this slice.
func ifftDIT(a []complex128) {
	twiddles := getTwiddles(len(a))

	for s := 1; s < bits.Len(uint(len(a))); s++ {
		m := 1 << s
		stride := len(a) / m
		half := m / 2

		for k := 0; k < len(a); k += m {
			for j := range half {
				twiddle := twiddles.inverse[j*stride]

				t := twiddle * a[k+j+half]
				u := a[k+j]

				a[k+j] = u + t
				a[k+j+half] = u - t
			}
		}
	}

	// Scale output by 1/n to return to original magnitude.
	invN := complex(1.0/float64(len(a)), 0)
	for i := range a {
		a[i] *= invN
	}
}

func getTwiddles(n int) twiddleTable {
	cacheMu.RLock()
	table, ok := cache[n]
	cacheMu.RUnlock()
	if ok {
		// Value is already stored in cache.
		return table
	}

	// We need to compute it.
	fwd := make([]complex128, n/2)
	inv := make([]complex128, n/2)
	for i := range n / 2 {
		theta := 2.0 * math.Pi * float64(i) / float64(n)
		s, c := math.Sincos(theta)
		fwd[i] = complex(c, -s)
		inv[i] = complex(c, s)
	}

	table = twiddleTable{forward: fwd, inverse: inv}

	go func() {
		cacheMu.Lock()
		cache[n] = table
		cacheMu.Unlock()
	}()

	return table
}
