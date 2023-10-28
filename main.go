package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/wcharczuk/go-chart"
)

func main() {
	arr := generateRandomArray(100000)

	// MergeSort
	inicio := time.Now()
	mergeSort(arr)
	tempoExecMerge := time.Since(inicio)
	fmt.Printf("Tempo de execução merge normal: %f\n", tempoExecMerge.Seconds())

	// MergeSortConcorrente (Abordagem 1)
	inicio = time.Now()
	mergeSortConcorrente(arr)
	tempoExecMergeConcorrente := time.Since(inicio)
	fmt.Printf("Tempo de execução merge concorrente (ineficiente): %f\n", tempoExecMergeConcorrente.Seconds())

	// MergeSortConcorrente (Abordagem 2 - Otimização de criação das goroutines)
	inicio = time.Now()
	mergeSortConcurrent(arr)
	tempoExecMergeConcorrente2 := time.Since(inicio)
	fmt.Printf("Tempo de execução merge concorrente (eficiente): %f\n", tempoExecMergeConcorrente2.Seconds())

	// QuickSort
	inicio = time.Now()
	quickSort(arr)
	tempoExecQuickSort := time.Since(inicio)
	fmt.Printf("Tempo de execução quickSort: %f\n", tempoExecQuickSort.Seconds())

	// QuickSortConcorrente (Abordagem 1)
	inicio = time.Now()
	quickSortConcorrente(arr)
	tempoExecQuickSortConcorrente := time.Since(inicio)
	fmt.Printf("Tempo de execução quickSort concorrente (ineficiente): %f\n", tempoExecQuickSortConcorrente.Seconds())

	// QuickSortConcorrente2 (Abordagem mais eficiente)
	inicio = time.Now()
	parallelQuickSort(arr, 0, len(arr)-1)
	tempoExecQuickSortConcorrente2 := time.Since(inicio)
	fmt.Printf("Tempo de execução quickSort concorrente (eficiente): %f\n", tempoExecQuickSortConcorrente2.Seconds())

	// Plot grafico de barras simples
	graph := chart.BarChart{
		Title: "Comparação dos algoritmos Merge e Quick Sort",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 100,
			},
		},
		Height:   512,
		BarWidth: 60,
		Bars: []chart.Value{
			{Value: float64(tempoExecMerge.Seconds()), Label: "Blue"},
			{Value: float64(tempoExecMergeConcorrente.Seconds()), Label: "Green"},
			{Value: float64(tempoExecMergeConcorrente2.Seconds()), Label: "Gray"},
			{Value: float64(tempoExecQuickSort.Seconds()), Label: "Orange"},
			{Value: float64(tempoExecQuickSortConcorrente.Seconds()), Label: "!"},
			{Value: float64(tempoExecQuickSortConcorrente2.Seconds()), Label: "Black"},
		},
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)

}

// ========== Merge Sort Tradicional
func mergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	middle := len(arr) / 2
	left := mergeSort(arr[:middle])
	right := mergeSort(arr[middle:])

	return merge(left, right)
}

// ========== Merge Sort com CSP e goroutines (Ineficiente)
func mergeSortConcorrente(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	middle := len(arr) / 2

	leftChan := make(chan []int)
	rightChan := make(chan []int)

	go func() {
		left := mergeSortConcorrente(arr[:middle])
		leftChan <- left
		close(leftChan)
	}()

	go func() {
		right := mergeSortConcorrente(arr[middle:])
		rightChan <- right
		close(rightChan)
	}()

	left, right := <-leftChan, <-rightChan

	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}

func generateRandomArray(size int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Intn(size * 10)
	}
	return arr
}

// ========== Merge Sort Eficiente (Otimização na criação de threads)
func mergeSortConcurrent(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	// Verifica o número de núcleos do processador
	numCores := runtime.NumCPU()
	runtime.GOMAXPROCS(numCores)

	// Divide o trabalho entre as goroutines
	chunks := chunkArray(arr, numCores)
	var wg sync.WaitGroup
	wg.Add(len(chunks))

	// Usa canais para coletar os resultados
	results := make([]chan []int, len(chunks))

	for i := 0; i < len(chunks); i++ {
		results[i] = make(chan []int)
		go func(i int) {
			defer wg.Done()
			results[i] <- mergeSort(chunks[i])
		}(i)
	}

	go func() {
		wg.Wait()
		for i := 0; i < len(chunks); i++ {
			close(results[i])
		}
	}()

	// Combina os resultados das goroutines
	sortedChunks := make([][]int, len(chunks))
	for i := 0; i < len(chunks); i++ {
		sortedChunks[i] = <-results[i]
	}

	return mergeSortedChunks(sortedChunks)
}

func chunkArray(arr []int, numChunks int) [][]int {
	chunkSize := (len(arr) + numChunks - 1) / numChunks
	chunks := make([][]int, numChunks)
	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		chunks[i] = arr[start:end]
	}
	return chunks
}

func mergeSortedChunks(chunks [][]int) []int {
	if len(chunks) == 0 {
		return []int{}
	}

	result := chunks[0]
	for i := 1; i < len(chunks); i++ {
		result = merge(result, chunks[i])
	}
	return result
}

// ========== Quick Sort Tradicional
func quickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	pivot := arr[len(arr)/2]
	left := make([]int, 0)
	right := make([]int, 0)
	equal := make([]int, 0)

	for _, value := range arr {
		if value < pivot {
			left = append(left, value)
		} else if value == pivot {
			equal = append(equal, value)
		} else {
			right = append(right, value)
		}
	}

	left = quickSort(left)
	right = quickSort(right)

	return append(append(left, equal...), right...)
}

// ========== Quick Sort Concorrente (Ineficiente)
func quickSortConcorrente(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	pivot := arr[len(arr)/2]
	left := make([]int, 0)
	right := make([]int, 0)
	equal := make([]int, 0)

	for _, value := range arr {
		if value < pivot {
			left = append(left, value)
		} else if value == pivot {
			equal = append(equal, value)
		} else {
			right = append(right, value)
		}
	}

	leftCh := make(chan []int)
	rightCh := make(chan []int)

	go func() {
		leftCh <- quickSortConcorrente(left)
	}()
	go func() {
		rightCh <- quickSortConcorrente(right)
	}()

	leftSorted, rightSorted := <-leftCh, <-rightCh
	return append(append(leftSorted, equal...), rightSorted...)
}

// ========== Quick Sort Concorrente (eficiente)
func parallelQuickSort(arr []int, low, high int) {
	if low < high {
		pivot := partition(arr, low, high)
		done := make(chan bool)

		go func() {
			parallelQuickSort(arr, low, pivot-1)
			done <- true
		}()
		parallelQuickSort(arr, pivot+1, high)

		<-done
	}
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return i
}
