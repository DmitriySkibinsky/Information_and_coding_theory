package main

import (
	"fmt"
	"math"
)

func entropyXI(p float64) float64 {
	entropy := math.Log2(1 / p)
	return entropy
}

func entropy(probabilities []float64) float64 {
	var h float64
	for _, p := range probabilities {
		if p > 0 {
			h += p * math.Log2(p)
		}
	}
	return -h
}

func informationContent(probabilities []float64) float64 {
	var info float64
	for _, p := range probabilities {
		if p > 0 {
			info += p * math.Log2(p)
		}
	}
	return -info
}

func main() {
	answer := entropyXI(0.0001)
	fmt.Println(answer)

}
