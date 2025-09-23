package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

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

func generateProbCorrect(n int, start, end float64) ([]float64, error) {
	if n <= 0 {
		return nil, fmt.Errorf("Число вероятностей безошибочной передачи должно быть больше 0")
	}
	rand.Seed(time.Now().UnixNano())
	probs := make([]float64, n)
	for i := 0; i < n; i++ {
		value := start + rand.Float64()*(end-start)
		probs[i] = value // Убрано округление
	}
	return probs, nil
}

// generateConditionalMatrix создает матрицу P(X|Y)
func generateConditionalMatrix(n int, probsRight []float64) ([][]float64, error) {
	if n <= 0 {
		return nil, fmt.Errorf("размер матрицы должен быть больше 0")
	}
	if len(probsRight) != n {
		return nil, fmt.Errorf("длина массива probsRight должна быть равна %d", n)
	}

	// Инициализация матрицы n x n
	matrix := make([][]float64, n)
	for i := range matrix {
		matrix[i] = make([]float64, n)
	}

	// Заполнение матрицы
	for i := 0; i < n; i++ {
		// Проверка, что вероятность безошибочной передачи в допустимом диапазоне
		if probsRight[i] < 0 || probsRight[i] > 1 {
			return nil, fmt.Errorf("вероятность probsRight[%d] = %f вне диапазона [0,1]", i, probsRight[i])
		}

		// Диагональный элемент P(X_i|Y_i)
		matrix[i][i] = probsRight[i]

		// Недиагональные элементы P(X_j|Y_i) для j != i
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

func calculateEntropy(n int, probs []float64) (float64, error) {
	if n <= 0 || len(probs) != n {
		return 0, fmt.Errorf("некорректные входные данные")
	}
	entropy := 0.0
	for i := 0; i < n; i++ {
		if probs[i] > 0 { // Избегаем log(0)
			entropy -= probs[i] * math.Log2(probs[i])
		}
	}
	return entropy, nil
}

func calculateConditionalEntropy(n int, jointMatrix [][]float64, outputProbs []float64) (float64, error) {
	if n <= 0 || len(jointMatrix) != n || len(jointMatrix[0]) != n || len(outputProbs) != n {
		return 0, fmt.Errorf("некорректные размеры данных")
	}

	entropy := 0.0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if jointMatrix[i][j] > 0 && outputProbs[j] > 0 { // Избегаем log(0) и деления на 0
				conditionalProb := jointMatrix[i][j] / outputProbs[j]
				entropy -= jointMatrix[i][j] * math.Log2(conditionalProb)
			}
		}
	}
	return entropy, nil
}

func runExperement(itr, n int) {
	for i := 0; i < itr; i++ {
		probs, err := generateProbabilities(n)
		if err != nil {
			panic(err)
		}
		fmt.Println(probs)

		probs_right, err := generateProbCorrect(n, 0.7, 1)
		if err != nil {
			panic(err)
		}
		matrix, _ := generateConditionalMatrix(n, probs_right)
		//for i := range res {
		//	fmt.Println(res[i])
		//}
		outputProbs, _ := calculateOutputProbabilities(n, probs, matrix)
		jointProbs, _ := calculateJointProbabilityMatrix(n, outputProbs, matrix) // как часто вместе появляются входные X и выходные Y сообщения
		fmt.Println(jointProbs)
		entropy, err := calculateEntropy(n, probs)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Энтропия на входе H(X) = %.4f бит\n", entropy)

		conditionalEntropy, err := calculateConditionalEntropy(n, jointProbs, outputProbs)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Условная энтропия H(X|Y) = %.4f бит\n", conditionalEntropy)

		fmt.Println("Количество информации при неполной достоверности сообщений ", entropy-conditionalEntropy)
	}
}

func main() {
	n := 53
	itr := 6
	runExperement(itr, n)
}
