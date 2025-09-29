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

// Вероятности достоверности сообщения
func generateProbCorrect(n int, start, end float64) ([]float64, error) {
	if n <= 0 {
		return nil, fmt.Errorf("Число вероятностей безошибочной передачи должно быть больше 0")
	}
	rand.Seed(time.Now().UnixNano())
	probs := make([]float64, n)
	for i := 0; i < n; i++ {
		value := start + rand.Float64()*(end-start)
		probs[i] = value
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

// Стадартная функция энтропии
func calculateEntropy(n int, probs []float64) (float64, error) {
	if n <= 0 || len(probs) != n {
		return 0, fmt.Errorf("некорректные входные данные")
	}
	entropy := 0.0
	for i := 0; i < n; i++ {
		if probs[i] > 0 {
			entropy -= probs[i] * math.Log2(probs[i])
		}
	}
	return entropy, nil
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

func runExperiment(itr, n int) {
	// Заголовок таблицы
	fmt.Println("+----------+---------------------+-------------------------------+------------------------------------+")
	fmt.Println("| Итерация | Энтропия H(X), бит  | Условная энтропия H(X|Y), бит | Количество информации I(X;Y), бит  |")
	fmt.Println("+----------+---------------------+-------------------------------+------------------------------------+")

	// Выполнение итераций и сбор данных
	for i := 0; i < itr; i++ {
		probs, err := generateProbabilities(n)
		if err != nil {
			panic(err)
		}

		probsRight, err := generateProbCorrect(n, 0.7, 1)
		if err != nil {
			panic(err)
		}

		matrix, err := generateConditionalMatrix(n, probsRight)
		if err != nil {
			panic(err)
		}

		outputProbs, err := calculateOutputProbabilities(n, probs, matrix)
		if err != nil {
			panic(err)
		}

		jointProbs, err := calculateJointProbabilityMatrix(n, outputProbs, matrix)
		if err != nil {
			panic(err)
		}

		entropy, err := calculateEntropy(n, probs)
		if err != nil {
			panic(err)
		}

		conditionalEntropy, err := calculateConditionalEntropy(n, jointProbs, outputProbs)
		if err != nil {
			panic(err)
		}

		// Форматирование строки таблицы
		fmt.Printf("| %8d | %19.4f | %29.4f | %34.4f |\n",
			i+1, entropy, conditionalEntropy, entropy-conditionalEntropy)
	}

	// Нижняя граница таблицы
	fmt.Println("+----------+---------------------+-------------------------------+------------------------------------+")
}

func main() {
	n := 53
	itr := 6
	runExperiment(itr, n)
}
