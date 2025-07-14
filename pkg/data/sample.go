package data

import (
	"log"
	"math/rand"
)

func Probabilities(counts []uint16) []float64 {
	sum := 0
	for _, c := range counts {
		sum += int(c)
	}

	ps := make([]float64, len(counts))
	for i, c := range counts {
		ps[i] = float64(c)/float64(sum)
	}

	return ps
}

func Sample(parameters []float64, generator *rand.Rand) int {
	X := generator.Float64()

	F := 0.0
	for i, f := range parameters {
		F += f

		if X <= F {
			return i
		}
	}

	log.Printf("WARN: Total probabilities %f < %f", F, X)

	return len(parameters) - 1
}
