package utils

import "math"

// FormatDecimal formata um float64 para ter precisão exata de 4 casas decimais
// nos cálculos, garantindo arredondamento adequado
func FormatDecimal(value float64) float64 {
	// Arredondamento mais preciso usando multiplicação por 10000 e depois round
	return math.Round(value*10000) / 10000
}

// FormatForDisplay formata um float64 para exibição com 2 casas decimais
// Esta função é usada apenas para saída/exibição, não para cálculos internos
func FormatForDisplay(value float64) float64 {
	return math.Round(value*100) / 100
}
