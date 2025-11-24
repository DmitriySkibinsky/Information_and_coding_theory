package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// === Глобальные константы и переменные ===
// k — количество информационных бит
const k = 4

// p — минимальное число проверочных бит, необходимое для кодирования k бит
var p = minP(k)

// n = 2^p - 1 — длина кода Хэмминга (всего бит)
var n = (1 << p) - 1

// nExt = n + 1 — длина расширенного кода (с общим паритетом)
var nExt = n + 1

// parityPos — позиции проверочных бит (степени двойки: 1, 2, 4, 8, ...)
var parityPos = pow2UpTo(n)

// dataPos — позиции информационных бит (все позиции от 1 до n, кроме проверочных)
var dataPos = makeDataPos()

// minP вычисляет минимальное p такое, что (2^p - p - 1) >= k
func minP(k int) int {
	p := 2
	for (1<<p)-p-1 < k {
		p++
	}
	return p
}

// pow2UpTo возвращает массив степеней двойки до max (включительно)
func pow2UpTo(max int) []int {
	var list []int
	for x := 1; x <= max; x <<= 1 {
		list = append(list, x)
	}
	return list
}

// makeDataPos формирует массив позиций информационных бит (все от 1 до n, кроме степеней 2)
func makeDataPos() []int {
	set := make(map[int]bool)
	for _, pos := range parityPos {
		set[pos] = true
	}
	var res []int
	for i := 1; i <= n; i++ {
		if !set[i] {
			res = append(res, i)
		}
	}
	return res
}

// randomBits генерирует случайный массив из len бит (0 или 1)
func randomBits(len int) []int {
	v := make([]int, len)
	for i := 0; i < len; i++ {
		v[i] = rand.Intn(2)
	}
	return v
}

// encodeSEC кодирует k информационных бит в код Хэмминга длины n
// Заполняет информационные позиции, затем вычисляет проверочные биты
func encodeSEC(dataK []int) []int {
	if len(dataK) != k {
		panic("Ожидались k бит данных")
	}
	word := make([]int, n)

	// Заполнение информационных бит
	for i := 0; i < k; i++ {
		word[dataPos[i]-1] = dataK[i]
	}

	// Вычисление проверочных бит
	for j := 0; j < p; j++ {
		parity := 0
		parityPosition := 1 << j
		for pos := 1; pos <= n; pos++ {
			if ((pos >> j) & 1) == 1 {
				parity ^= word[pos-1]
			}
		}
		word[parityPosition-1] = parity & 1
	}
	return word
}

// addOverallParity добавляет общий паритетный бит к коду Хэмминга (SEC → SECDED)
func addOverallParity(hammingWord []int) []int {
	overall := 0
	for _, v := range hammingWord {
		overall ^= v
	}
	return append(hammingWord, overall)
}

// extractData извлекает k информационных бит из исправленного кодового слова длины n
func extractData(hammingWordN []int) []int {
	data := make([]int, k)
	for i := 0; i < k; i++ {
		data[i] = hammingWordN[dataPos[i]-1]
	}
	return data
}

// buildH строит проверочную матрицу H размера p x n для кода Хэмминга
// H[i][j] = 1, если (j+1)-й бит участвует в i-м проверочном уравнении
func buildH(n, p int) [][]int {
	H := newMatrix(p, n)
	for col := 1; col <= n; col++ {
		for row := 0; row < p; row++ {
			H[row][col-1] = (col >> row) & 1
		}
	}
	return H
}

// multiply выполняет умножение матрицы H на вектор vec по модулю 2 (XOR)
func multiply(H [][]int, vec []int) []int {
	rows, cols := len(H), len(H[0])
	s := make([]int, rows)
	for i := 0; i < rows; i++ {
		acc := 0
		for j := 0; j < cols; j++ {
			acc ^= (H[i][j] & vec[j])
		}
		s[i] = acc & 1
	}
	return s
}

// syndromeIndex преобразует синдром (массив бит) в десятичное число — позицию ошибки
func syndromeIndex(s []int) int {
	val := 0
	for i, bit := range s {
		if bit == 1 {
			val |= 1 << i
		}
	}
	return val
}

// newMatrix создаёт двумерный срез (матрицу) размера r x c, заполненный нулями
func newMatrix(r, c int) [][]int {
	m := make([][]int, r)
	for i := range m {
		m[i] = make([]int, c)
	}
	return m
}

// bitsToString преобразует массив бит в строку из 0 и 1
func bitsToString(bits []int) string {
	var sb strings.Builder
	for _, b := range bits {
		sb.WriteString(fmt.Sprintf("%d", b))
	}
	return sb.String()
}

// markErrorPositions создаёт строку с метками ^ под позициями ошибок
func markErrorPositions(length int, errorPositions []int) string {
	marks := make([]byte, length)
	for i := range marks {
		marks[i] = ' '
	}
	for _, pos := range errorPositions {
		if pos-1 < length {
			marks[pos-1] = '^'
		}
	}
	result := string(marks)
	result += " (позиции "
	posStrs := make([]string, len(errorPositions))
	for i, pos := range errorPositions {
		posStrs[i] = fmt.Sprintf("%d", pos)
	}
	result += strings.Join(posStrs, ", ")
	result += ")"
	return result
}

// printMatrix выводит матрицу с именем, ограничивая вывод 15 столбцами
func printMatrix(M [][]int, name string) {
	fmt.Printf("\nМатрица %s:\n", name)
	rows, cols := len(M), len(M[0])

	maxColsToShow := 15
	if cols > maxColsToShow {
		fmt.Printf("(Отображаются первые %d из %d столбцов)\n\n", maxColsToShow, cols)
	}

	for i := 0; i < rows; i++ {
		limit := int(math.Min(float64(cols), float64(maxColsToShow)))
		for j := 0; j < limit; j++ {
			fmt.Printf("%d ", M[i][j])
		}
		if cols > maxColsToShow {
			fmt.Print("...")
		}
		fmt.Println()
	}
}

// bitsEqual сравнивает два массива бит на полное совпадение
func bitsEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// === Основная функция программы ===
func main() {
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Вывод начальной информации о параметрах кода
	fmt.Printf("k = %d, подобрано p = %d, Длина кода Хэмминга = %d\n", k, p, n)
	fmt.Println("SEC:  исправление 1 ошибки (d_min = 3)")
	fmt.Println("SECDED: обнаружение 2-х ошибок (d_min = 4)\n")

	// Построение проверочной матрицы H для кода Хэмминга
	H := buildH(n, p)
	fmt.Printf("Проверочная матрица H: %d x %d\n", p, n)

	// Количество экспериментов
	experiments := 10
	for t := 1; t <= experiments; t++ {
		fmt.Println("")
		fmt.Println("")
		fmt.Println(strings.Repeat("_", 70))
		fmt.Printf("\nЭксперимент #%d\n", t)

		// 1. Генерация случайных данных длиной k бит
		a := randomBits(k)
		fmt.Println("\nИсходная кодовая комбинация", strings.Join(toStrings(a), ""))

		// 2. Кодирование данных в код Хэмминга (SEC)
		cSec := encodeSEC(a)
		fmt.Println("Закодированная кодовая комбинация", strings.Join(toStrings(cSec), ""))

		// 3. Добавление общего паритетного бита → SECDED
		cSecded := addOverallParity(cSec)

		// 4. Генерация 0, 1 или 2 ошибок в кодовом слове
		multiplicity := rand.Intn(3)
		r := append([]int(nil), cSecded...)
		errorPositions := []int{}

		if multiplicity == 1 {
			pos := rand.Intn(nExt) + 1
			r[pos-1] ^= 1
			errorPositions = append(errorPositions, pos)
		} else if multiplicity == 2 {
			p1 := rand.Intn(nExt) + 1
			p2 := 0
			for p2 == 0 || p2 == p1 {
				p2 = rand.Intn(nExt) + 1
			}
			r[p1-1] ^= 1
			r[p2-1] ^= 1
			errorPositions = append(errorPositions, p1, p2)
			sort.Ints(errorPositions)
		}

		// Вывод информации о внесённых ошибках
		errStr := "-"
		if len(errorPositions) > 0 {
			posStrs := make([]string, len(errorPositions))
			for i, pos := range errorPositions {
				posStrs[i] = fmt.Sprintf("%d", pos)
			}
			errStr = strings.Join(posStrs, ",")
		}
		fmt.Printf("\nВнесение ошибки: %d (позиции: %s)\n", multiplicity, errStr)
		fmt.Printf("Информационная комбинация с ошибкой: %s\n", strings.Join(toStrings(r), ""))

		// 5. Разделение: общий паритет и кодовое слово SEC
		//
		//overall := r[nExt-1]
		rSec := r[:n]

		// 6. Вычисление синдрома: S = H * rSec (по модулю 2)
		syndrome := multiply(H, rSec)

		// 7. Преобразование синдрома в индекс ошибки
		synIndex := syndromeIndex(syndrome)

		// 8. Вычисление общей чётности принятого слова
		overallParity := 0
		for _, v := range r {
			overallParity ^= v
		}

		// 9. Анализ синдрома и общей чётности → принятие решения
		var verdict string
		correctedPos := 0
		rFixed := append([]int(nil), r...)

		if synIndex == 0 && overallParity == 0 {
			verdict = "Ошибок нет"
		} else if synIndex != 0 && overallParity == 1 {
			verdict = fmt.Sprintf("Одиночная ошибка (позиция %d)", synIndex)
			rFixed[synIndex-1] ^= 1
			correctedPos = synIndex
		} else if synIndex == 0 && overallParity == 1 {
			verdict = fmt.Sprintf("Ошибка только в общем паритетном бите (позиция %d)", nExt)
			rFixed[nExt-1] ^= 1
			correctedPos = nExt
		} else {
			verdict = "Двукратная ошибка (обнаружена, исправить нельзя)"
		}

		// Вывод синдрома и решения
		synStr := ""
		for i := len(syndrome) - 1; i >= 0; i-- {
			synStr += fmt.Sprintf("%d", syndrome[i])
		}
		fmt.Printf("Синдром (SEC): %s -> %d\n", synStr, synIndex)
		printH()

		fmt.Printf("\nРешение: %s\n", verdict)

		if strings.HasPrefix(verdict, "Одиночная") || strings.HasPrefix(verdict, "Ошибка только") {
			aDecoded := extractData(rFixed[:n])
			ok := bitsEqual(a, aDecoded)

			fmt.Printf("Данные восстановлены корректно: %t  (исправлена позиция %d)\n\n", ok, correctedPos)

			if multiplicity == 1 || correctedPos == nExt {
				// Красивый вывод трёх строк: оригинал → принятое → исправленное
				orig := bitsToSpacedString(cSecded)
				recv := bitsToSpacedString(r)
				fixed := bitsToSpacedString(rFixed)

				fmt.Println("Визуализация исправления ошибки:")
				fmt.Printf("Исходное кодовое слово : %s\n", orig)
				fmt.Printf("Принятое (с ошибкой)   : %s\n", recv)
				fmt.Printf("Исправленное слово     : %s\n", fixed)

				// Стрелка под ошибкой
				arrow := strings.Repeat(" ", (correctedPos-1)*2) + "↑"
				fmt.Printf("                         %s (позиция %d)\n", arrow, correctedPos)

				// Если ошибка была в общем паритете — отдельно укажем
				if correctedPos == nExt {
					fmt.Printf("                           ↑ ошибка в общем паритетном бите\n")
				}
			}

			fmt.Println()

		} else if verdict == "Двукратная ошибка (обнаружена, исправить нельзя)" {
			// 11. Визуализация двойной ошибки
			fmt.Printf("Исходная:  %s\n", bitsToString(cSecded))
			fmt.Printf("Принятая:  %s\n", bitsToString(r))
			fmt.Printf("Ошибки:    %s\n", markErrorPositions(len(cSecded), errorPositions))
		}
	}
}

// === Вспомогательные функции ===
func toStrings(arr []int) []string {
	res := make([]string, len(arr))
	for i, v := range arr {
		res[i] = fmt.Sprintf("%d", v)
	}
	return res
}

// Превращает []int в строку вида "1 0 1 1 0 0 1 0 1" — с пробелами
func bitsToSpacedString(bits []int) string {
	strs := make([]string, len(bits))
	for i, b := range bits {
		strs[i] = fmt.Sprintf("%d", b)
	}
	return strings.Join(strs, " ")
}
func printH() {
	H := buildH(n, p)
	fmt.Printf("\n             ПРОВЕРОЧНАЯ МАТРИЦА H  (%d × %d)\n", p, n)
	fmt.Print("Позиции → ")
	for i := 1; i <= n; i++ {
		fmt.Printf("%2d ", i)
	}
	fmt.Println("\n          " + strings.Repeat("─", n*3))

	for i := 0; i < p; i++ {
		pb := 1 << i
		fmt.Printf("  P%-2d →  ", pb)
		for j := 0; j < n; j++ {
			if H[i][j] == 1 {
				fmt.Print(" 1 ")
			} else {
				fmt.Print(" 0 ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
