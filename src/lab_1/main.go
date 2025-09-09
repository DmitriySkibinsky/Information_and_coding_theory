package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"gonum.org/v1/gonum/stat"
)

func entropy(probabilities []float64) float64 {
	var h float64
	for _, p := range probabilities {
		if p > 0 {
			h += p * math.Log2(p)
		}
	}
	return -h
}

func generateProbabilities(n int) ([]float64, error) {
	if n <= 0 {
		return nil, fmt.Errorf("число вероятностей должно быть больше 0")
	}
	rand.Seed(time.Now().UnixNano())
	probs := make([]float64, n)
	sum := 0.0
	for i := 0; i < n; i++ {
		probs[i] = rand.Float64()
		sum += probs[i]
	}
	for i := 0; i < n; i++ {
		probs[i] /= sum
	}
	return probs, nil
}

func maxEntropy(n int) float64 {
	return math.Log2(float64(n))
}

func runExperiment(n int, experimentNum int) (float64, float64, error) {
	probs, err := generateProbabilities(n)
	if err != nil {
		fmt.Printf("Эксперимент %d: Ошибка: %v\n", experimentNum, err)
		return 0, 0, err
	}
	avgEntropy := entropy(probs)
	maxEnt := maxEntropy(n)

	// Преобразуем массив вероятностей в строку с округлением
	probStrs := make([]string, len(probs))
	for i, p := range probs {
		probStrs[i] = fmt.Sprintf("%.4f", p)
	}
	probJoined := "[" + strings.Join(probStrs, ", ") + "]"

	// Табличный вывод
	fmt.Printf("| %3d | %3d | %-110s | %10.4f | %10.4f |\n",
		experimentNum, n, probJoined, avgEntropy, maxEnt)

	return avgEntropy, maxEnt, nil
}

func main() {
	// Заголовок таблицы
	fmt.Println("| Exp |  n  | Вероятности                                                                                                    |  Средн. H  |  Макс. H   |")
	fmt.Println("|-----|-----|----------------------------------------------------------------------------------------------------------------|------------|------------|")

	ns := []int{8, 9, 10, 11, 12, 13}
	avgAvgEntropy := []float64{}
	avgMaxEnt := []float64{}

	for i, n := range ns {
		avgEntropy, maxEnt, err := runExperiment(n, i+1)
		if err == nil {
			avgAvgEntropy = append(avgAvgEntropy, avgEntropy)
			avgMaxEnt = append(avgMaxEnt, maxEnt)
		}
	}

	// Вывод списков
	fmt.Println("\nСписок Средн. H:", stat.Mean(avgAvgEntropy, nil))

}
