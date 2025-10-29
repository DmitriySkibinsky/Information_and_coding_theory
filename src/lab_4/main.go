package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func calculateN(k int) int {
	n := k + 1
	for math.Pow(2, float64(k)) > math.Pow(2, float64(n))/(1+float64(n)) {
		n++
	}
	return n
}

// === 1. Генерация случайного информационного сообщения длины k ===
func generateMessage(k int) []int {
	msg := make([]int, k)
	for i := 0; i < k; i++ {
		msg[i] = rand.Intn(2)
	}
	return msg
}

// === 2.1 Построение производящей матрицы H_p,k = [U_k | H_p] ===
func buildGeneratorMatrix(n, p, k int) [][]int {
	G := make([][]int, k)
	for i := range G {
		G[i] = make([]int, n)
	}

	d := 3
	for i := 0; i < k; i++ {
		// U_k — единичная матрица
		for j := 0; j < k; j++ {
			G[i][j] = boolToInt(i == j)
		}

		// Увеличиваем d при достижении степени 2
		if i+d == int(math.Pow(2, float64(d-1))) {
			d++
		}

		if i < k-1 {
			num := i + d
			binary := fmt.Sprintf("%0*b", p, num)

			for j := 0; j < p; j++ {
				G[i][k+j] = int(binary[j] - '0')
			}
		} else {
			for j := 0; j < p; j++ {
				G[i][k+j] = 1
			}
		}
	}
	return G
}

// === 2.2 Построение проверочной матрицы H = [H_p | I_p] ===
func buildParityCheckMatrix(G [][]int, n, p, k int) [][]int {
	H := make([][]int, p)
	for i := range H {
		H[i] = make([]int, n)
	}

	// Копируем H_p
	for i := 0; i < p; i++ {
		for j := 0; j < k; j++ {
			H[i][j] = G[j][k+i]
		}
	}

	// I_p — единичная матрица в правой части
	for i := 0; i < p; i++ {
		for j := 0; j < p; j++ {
			H[i][k+j] = boolToInt(i == j)
		}
	}
	return H
}

// === 2.3 Кодирование: добавление проверочных битов ===
func encodeMessage(msg []int, H [][]int, n, p, k int) []int {
	codeword := make([]int, n)
	copy(codeword, msg)

	// Вычисляем проверочные биты
	for i := 0; i < p; i++ {
		sum := 0
		//printVector("code ", codeword)
		for j := 0; j < k; j++ {

			if H[i][j] == 1 {
				sum += codeword[j]
			}
		}
		codeword[k+i] = sum % 2
	}
	return codeword
}

// === 4. Внесение однократной ошибки (только в информационной части) ===
func injectError(codeword []int, k int) (int, []int) {
	errPos := rand.Intn(k)
	noisy := make([]int, len(codeword))
	copy(noisy, codeword)
	noisy[errPos] = 1 - noisy[errPos]
	return errPos, noisy
}

// === 5. Вычисление синдрома ===
func computeSyndrome(received []int, H [][]int, n, p, k int) []int {
	syndrome := make([]int, p)
	for i := 0; i < p; i++ {
		sum := received[k+i]
		for j := 0; j < k; j++ {
			if H[i][j] == 1 {
				sum += received[j]
			}
		}
		syndrome[i] = sum % 2
	}
	return syndrome
}

// === 6. Поиск позиции ошибки по синдрому ===
func findErrorPosition(syndrome []int, H [][]int, p, k int) (int, bool) {
	for pos := 0; pos < k; pos++ {
		match := true
		for j := 0; j < p; j++ {
			if H[j][pos] != syndrome[j] {
				match = false
				break
			}
		}
		if match {
			return pos, true
		}
	}
	return -1, false
}

// === Утилита: bool → int ===
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// === Вывод вектора ===
func printVector(label string, v []int) {
	fmt.Printf("%s", label)
	for _, x := range v {
		fmt.Printf("%d ", x)
	}
	fmt.Println()
}

// === Основная функция эксперимента ===
func runExperiment(exp, k int) {
	fmt.Printf("\n=== Эксперимент %d ===\n", exp+1)

	n := calculateN(k)
	p := n - k

	// 1. Генерация сообщения
	infoMsg := generateMessage(k)
	fmt.Printf("Информационное сообщение (k=%d): ", k)
	printVector("", infoMsg)

	// 2. Построение матриц
	G := buildGeneratorMatrix(n, p, k)
	//for i := range G {
	//	fmt.Println(G[i])
	//}
	fmt.Println()
	H := buildParityCheckMatrix(G, n, p, k)
	//for i := range H {
	//	fmt.Println(H[i])
	//}

	// 3. Кодирование
	codeword := encodeMessage(infoMsg, H, n, p, k)
	printVector("Кодовое слово (n бит):     ", codeword)

	// 4. Внесение ошибки
	errPos, noisy := injectError(codeword, k)
	printVector("С ошибкой:                 ", noisy)
	fmt.Printf("Ошибка в позиции: %d\n", errPos)

	// 5. Синдром
	syndrome := computeSyndrome(noisy, H, n, p, k)
	printVector("Синдром:                   ", syndrome)

	// 6. Обнаружение и исправление
	foundPos, found := findErrorPosition(syndrome, H, p, k)
	if found {
		fmt.Printf("Ошибка обнаружена и исправлена в позиции: %d\n", foundPos)
		// Исправляем
		corrected := make([]int, n)
		copy(corrected, noisy)
		corrected[foundPos] = 1 - corrected[foundPos]
		printVector("Исправленное сообщение:    ", corrected)
	} else {
		fmt.Println("Ошибка не обнаружена (возможно, в проверочном бите)")
	}

	fmt.Println("----------------------------------------")
}

func main() {
	k := 53
	rand.Seed(time.Now().UnixNano())

	fmt.Printf("Код Хэмминга: k = %d информационных битов\n", k)
	fmt.Println("Запуск 8 экспериментов с обнаружением и исправлением однократных ошибок...\n")

	for exp := 0; exp < 8; exp++ {
		runExperiment(exp, k)
	}
}
