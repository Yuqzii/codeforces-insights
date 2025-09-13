package fft

func FloatToComplex(x []float64) []complex128 {
	res := make([]complex128, 0, len(x))
	for _, v := range x {
		res = append(res, complex(v, 0))
	}
	return res
}

func ComplexToFloat(x []complex128) []float64 {
	res := make([]float64, 0, len(x))
	for _, v := range x {
		res = append(res, real(v))
	}
	return res
}
