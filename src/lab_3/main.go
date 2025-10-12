package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Вероятности дискретных сообщений
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

func CalculateEntropy(n int, probs []float64) (float64, error) {
	if len(probs) == 0 || n <= 0 {
		return 0.0, fmt.Errorf("ошибка в входных данных")
	}
	entropy := 0.0
	for i := 0; i < n; i++ {
		if probs[i] > 0 {
			entropy -= probs[i] * math.Log2(probs[i])
		}
	}
	return entropy, nil
}

func generateDuration(n int, start, end float64) ([]float64, error) {
	if n <= 0 {
		return nil, fmt.Errorf("число сообщений для расчета длительности должно быть больше 0")
	}
	rand.Seed(time.Now().UnixNano())
	probs := make([]float64, n)
	for i := 0; i < n; i++ {
		value := start + rand.Float64()*(end-start)
		probs[i] = value
	}
	return probs, nil
}

func generateMiddleDuration(n int, probs []float64, massiveDuraions []float64) (float64, error) {
	if n <= 0 {
		return 0.0, fmt.Errorf("число сообщений для расчета средней длительности должно быть больше 0")
	}
	middleDuration := 0.0
	for i := 0; i < n; i++ {
		middleDuration += probs[i] * massiveDuraions[i]
	}
	return middleDuration, nil
}

// Вероятности достоверности сообщения
func generateProbCorrect(n int, start, end float64) ([]float64, error) {
	if n <= 0 {
		return nil, fmt.Errorf("число вероятностей безошибочной передачи должно быть больше 0")
	}
	rand.Seed(time.Now().UnixNano())
	probs := make([]float64, n)
	for i := 0; i < n; i++ {
		value := start + rand.Float64()*(end-start)
		probs[i] = value
	}
	return probs, nil
}

// Вероятности без помех (все 1.0)
func generateProbCorrectNoNoise(n int) ([]float64, error) {
	if n <= 0 {
		return nil, fmt.Errorf("число вероятностей безошибочной передачи должно быть больше 0")
	}
	probs := make([]float64, n)
	for i := 0; i < n; i++ {
		probs[i] = 1.0
	}
	return probs, nil
}

// Матрица вероятностей - если на входе Xi а на выходе Yj
func generateConditionalMatrix(n int, probsRight []float64) ([][]float64, error) {
	if n <= 0 {
		return nil, fmt.Errorf("размер матрицы должен быть больше 0")
	}
	if len(probsRight) != n {
		return nil, fmt.Errorf("длина массива probsRight должна быть равна %d", n)
	}

	matrix := make([][]float64, n)
	for i := range matrix {
		matrix[i] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		if probsRight[i] < 0 || probsRight[i] > 1 {
			return nil, fmt.Errorf("вероятность probsRight[%d] = %f вне диапазона [0,1]", i, probsRight[i])
		}
		matrix[i][i] = probsRight[i]
		if n > 1 {
			remainingProb := (1.0 - probsRight[i]) / float64(n-1)
			for j := 0; j < n; j++ {
				if j != i {
					matrix[i][j] = remainingProb
				}
			}
		}
	}
	return matrix, nil
}

// Вероятности появления выходных символов Xi с учётом возможных ошибок
func calculateOutputProbabilities(n int, inputProbs []float64, condMatrix [][]float64) ([]float64, error) {
	if n <= 0 || len(inputProbs) != n || len(condMatrix) != n || len(condMatrix[0]) != n {
		return nil, fmt.Errorf("некорректные размеры данных")
	}

	outputProbs := make([]float64, n)
	for j := 0; j < n; j++ {
		sum := 0.0
		for i := 0; i < n; i++ {
			sum += inputProbs[i] * condMatrix[i][j]
		}
		outputProbs[j] = sum
	}
	return outputProbs, nil
}

func calculateBandwidthCapacity(n int, middleDuration float64, conditionEntropy float64) (float64, error) {
	return math.Log2(float64(n)-conditionEntropy) / middleDuration, nil
}

// Это таблица, которая показывает, как часто конкретная пара "входной символ + выходной символ" встречается вместе.
func calculateJointProbabilityMatrix(n int, outputProbs []float64, condMatrix [][]float64) ([][]float64, error) {
	if n <= 0 || len(outputProbs) != n || len(condMatrix) != n || len(condMatrix[0]) != n {
		return nil, fmt.Errorf("некорректные размеры данных")
	}

	jointMatrix := make([][]float64, n)
	for i := range jointMatrix {
		jointMatrix[i] = make([]float64, n)
		for j := range jointMatrix[i] {
			jointMatrix[i][j] = outputProbs[j] * condMatrix[i][j]
		}
	}
	return jointMatrix, nil
}

// Условная энтропия выходного сообщения
func calculateConditionalEntropy(n int, jointMatrix [][]float64, outputProbs []float64) (float64, error) {
	if n <= 0 || len(jointMatrix) != n || len(jointMatrix[0]) != n || len(outputProbs) != n {
		return 0, fmt.Errorf("некорректные размеры данных")
	}

	entropy := 0.0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if jointMatrix[i][j] > 0 && outputProbs[j] > 0 {
				conditionalProb := jointMatrix[i][j] / outputProbs[j]
				entropy -= jointMatrix[i][j] * math.Log2(conditionalProb)
			}
		}
	}
	return entropy, nil
}

func calculateBaudRate(entropy float64, conditionalEntropy float64, middleDuration float64) (float64, error) {
	return (entropy - conditionalEntropy) / middleDuration, nil
}

func RunTests(n int) string {
	// Структура для хранения результатов
	type Result struct {
		Entropy            float64
		ConditionalEntropy float64
		MiddleDuration     float64
		BandwidthCapacity  float64
		BaudRate           float64
	}

	// Массивы для хранения результатов
	resultsWithNoise := make([]Result, 6)
	resultsNoNoise := make([]Result, 6)
	q := 1 - (1 / float64(n*2))

	// Тест с помехами
	for i := 0; i < 6; i++ {
		probs, _ := generateProbabilities(n)
		entropy, _ := CalculateEntropy(n, probs)
		massiveDuration, _ := generateDuration(n, 0, float64(n))
		probsRight, _ := generateProbCorrect(n, 0, q)
		matrix, _ := generateConditionalMatrix(n, probsRight)
		outputProbs, _ := calculateOutputProbabilities(n, probs, matrix)
		jointProbs, _ := calculateJointProbabilityMatrix(n, outputProbs, matrix)
		conditionEntropy, _ := calculateConditionalEntropy(n, jointProbs, outputProbs)
		middleDuration, _ := generateMiddleDuration(n, probs, massiveDuration)
		bandwidthCapacity, _ := calculateBandwidthCapacity(n, middleDuration, conditionEntropy)
		baudRate, _ := calculateBaudRate(entropy, conditionEntropy, middleDuration)

		resultsWithNoise[i] = Result{
			Entropy:            entropy,
			ConditionalEntropy: conditionEntropy,
			MiddleDuration:     middleDuration,
			BandwidthCapacity:  bandwidthCapacity,
			BaudRate:           baudRate,
		}
	}

	// Тест без помех
	for i := 0; i < 6; i++ {
		probs, _ := generateProbabilities(n)
		entropy, _ := CalculateEntropy(n, probs)
		massiveDuration, _ := generateDuration(n, 0, float64(n))
		probsRight, _ := generateProbCorrectNoNoise(n)
		matrix, _ := generateConditionalMatrix(n, probsRight)
		outputProbs, _ := calculateOutputProbabilities(n, probs, matrix)
		jointProbs, _ := calculateJointProbabilityMatrix(n, outputProbs, matrix)
		conditionEntropy, _ := calculateConditionalEntropy(n, jointProbs, outputProbs)
		middleDuration, _ := generateMiddleDuration(n, probs, massiveDuration)
		bandwidthCapacity, _ := calculateBandwidthCapacity(n, middleDuration, conditionEntropy)
		baudRate, _ := calculateBaudRate(entropy, conditionEntropy, middleDuration)

		resultsNoNoise[i] = Result{
			Entropy:            entropy,
			ConditionalEntropy: conditionEntropy,
			MiddleDuration:     middleDuration,
			BandwidthCapacity:  bandwidthCapacity,
			BaudRate:           baudRate,
		}
	}

	// Вычисление средних значений для канала с помехами
	var avgBandwidthWithNoise, avgBaudRateWithNoise float64
	for _, res := range resultsWithNoise {
		avgBandwidthWithNoise += res.BandwidthCapacity
		avgBaudRateWithNoise += res.BaudRate
	}
	avgBandwidthWithNoise /= 6.0
	avgBaudRateWithNoise /= 6.0

	// Вычисление средних значений для канала без помех
	var avgBandwidthNoNoise, avgBaudRateNoNoise float64
	for _, res := range resultsNoNoise {
		avgBandwidthNoNoise += res.BandwidthCapacity
		avgBaudRateNoNoise += res.BaudRate
	}
	avgBandwidthNoNoise /= 6.0
	avgBaudRateNoNoise /= 6.0

	// Формирование таблиц в формате Markdown
	resultStr := "# Результаты тестов\n\n"
	resultStr += fmt.Sprintf("## Тест с помехами (вероятность безошибочной передачи: [0, %.5f])\n", q)
	resultStr += "| Эксперимент | Энтропия H(X) | Условная энтропия H(Y|X) | Средняя длительность T (с) | Пропускная способность C (бит/с) | Скорость передачи R (бит/с) |\n"
	resultStr += "|-------------|---------------|--------------------------|---------------------------|----------------------------------|------------------------------|\n"
	for i, res := range resultsWithNoise {
		resultStr += fmt.Sprintf("| %11d | %13.4f | %24.4f | %25.4f | %32.4f | %28.4f |\n",
			i+1, res.Entropy, res.ConditionalEntropy, res.MiddleDuration, res.BandwidthCapacity, res.BaudRate)
	}
	resultStr += fmt.Sprintf("\n**Средняя пропускная способность C (бит/с):** %.4f\n", avgBandwidthWithNoise)
	resultStr += fmt.Sprintf("**Средняя скорость передачи R (бит/с):** %.4f\n", avgBaudRateWithNoise)

	resultStr += "\n## Тест без помех (вероятность безошибочной передачи: 1.0)\n"
	resultStr += "| Эксперимент | Энтропия H(X) | Условная энтропия H(Y|X) | Средняя длительность T (с) | Пропускная способность C (бит/с) | Скорость передачи R (бит/с) |\n"
	resultStr += "|-------------|---------------|--------------------------|---------------------------|----------------------------------|------------------------------|\n"
	for i, res := range resultsNoNoise {
		resultStr += fmt.Sprintf("| %11d | %13.4f | %24.4f | %25.4f | %32.4f | %28.4f |\n",
			i+1, res.Entropy, res.ConditionalEntropy, res.MiddleDuration, res.BandwidthCapacity, res.BaudRate)
	}
	resultStr += fmt.Sprintf("\n**Средняя пропускная способность C (бит/с):** %.4f\n", avgBandwidthNoNoise)
	resultStr += fmt.Sprintf("**Средняя скорость передачи R (бит/с):** %.4f\n", avgBaudRateNoNoise)

	return resultStr
}

func main() {
	n := 16
	result := RunTests(n)
	fmt.Println(result)
}
