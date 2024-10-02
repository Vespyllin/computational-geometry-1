package main

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
	"time"
)

func formatWithCommas(value float64) string {
	parts := strings.Split(fmt.Sprintf("%.2f", value), ".") // Split the float into integer and decimal parts
	integerPart := parts[0]
	decimalPart := parts[1]

	// Add commas to the integer part
	var result strings.Builder
	for i, v := range integerPart {
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(v)
	}

	// Combine integer part with decimal part
	return result.String() + "." + decimalPart
}

func genRandUniqueArr(size int) []int {
	numSet := make(map[int]struct{})
	result := make([]int, 0, size)

	for len(result) < size {
		n := rand.Int()
		if _, exists := numSet[n]; !exists {
			numSet[n] = struct{}{}
			result = append(result, n)
		}
	}
	return result
}

func tabulateResults(sizes []int, arrays [][][]int) {
	// SIZE
	for j := 0; j < len(sizes); j++ {
		fmt.Printf("SIZE: %d\t", sizes[j])
		// ALGO
		for i := 0; i < len(arrays); i++ {
			sum := 0

			// ITER
			for k := 0; k < len(arrays[i][j]); k++ {
				sum += arrays[i][j][k]
			}

			average := float64(sum) / float64(len(arrays[i][j]))
			fmt.Printf("%16s\t", formatWithCommas(average))
		}
		fmt.Println()
	}
}

func basicSortBenchmark(sizes []int, iter int) [][]int {
	basicSortDataNs := make([][]int, len(sizes))
	for i := range sizes {
		basicSortDataNs[i] = make([]int, iter)
	}

	for s, size := range sizes {
		for i := 0; i < iter; i++ {
			input := genRandUniqueArr(size)
			scratch := make([]int, len(input))

			startTime := time.Now()
			_, bmRres := basicMergeSort(input, scratch)
			elapsedTime := time.Since(startTime)

			if !sort.IntsAreSorted(bmRres) {
				panic(fmt.Sprintf("ERROR IN BASIC MERGE SORT\n\tSIZE: %d - ITERATION: %d\n", size, i))
			}

			basicTimeNs := elapsedTime.Nanoseconds()
			basicSortDataNs[s][i] = int(basicTimeNs)
		}
	}

	return basicSortDataNs
}

func parallelSortBenchmark(sizes []int, iter int, threads int) [][]int {

	parallelSortDataNs := make([][]int, len(sizes))
	for i := range sizes {
		parallelSortDataNs[i] = make([]int, iter)
	}

	for s, size := range sizes {
		for i := 0; i < iter; i++ {
			input := genRandUniqueArr(size)
			scratch := make([]int, len(input))

			startTime := time.Now()
			_, pmRres := segmentMergeSort(input, scratch, threads)
			elapsedTime := time.Since(startTime)

			if !sort.IntsAreSorted(pmRres) {
				panic(fmt.Sprintf("ERROR IN BASIC MERGE SORT\n\tSIZE: %d - ITERATION: %d\n", size, i))
			}

			parallelTimeNs := elapsedTime.Nanoseconds()
			parallelSortDataNs[s][i] = int(parallelTimeNs)
		}
	}

	return parallelSortDataNs
}

func sequentialMergeBenchmark(sizes []int, iter int) [][]int {
	mergeDataNs := make([][]int, len(sizes))
	for i := range sizes {
		mergeDataNs[i] = make([]int, iter)
	}

	for s, size := range sizes {
		for i := 0; i < iter; i++ {
			randomArr := genRandUniqueArr(size)
			sortedA := randomArr[len(randomArr)/2:]
			sortedB := randomArr[:len(randomArr)/2]
			sort.Ints(sortedA)
			sort.Ints(sortedB)

			dest := make([]int, len(randomArr))

			startTime := time.Now()
			sequentialMerge(sortedA, sortedB, dest)
			elapsedTime := time.Since(startTime)

			seqTimeNs := elapsedTime.Nanoseconds()
			mergeDataNs[s][i] = int(seqTimeNs)
		}
	}

	return mergeDataNs
}

func segmentMergeBenchmark(sizes []int, iter int, threads int) [][]int {
	mergeDataNs := make([][]int, len(sizes))
	for i := range sizes {
		mergeDataNs[i] = make([]int, iter)
	}

	for s, size := range sizes {
		for i := 0; i < iter; i++ {
			randomArr := genRandUniqueArr(size)
			inputA := randomArr[len(randomArr)/2:]
			inputB := randomArr[:len(randomArr)/2]
			sort.Ints(inputA)
			sort.Ints(inputB)

			rB := rank(inputB[(len(inputB)-1)*3/4], inputA)
			rA := rank(inputA[(len(inputA)-1)*3/4], inputB)
			if rB < len(inputA)/2 || rA < len(inputB)/2 {
				fmt.Printf("Bad distribution at iteration: %d of size %d", iter, size)
				i--
				continue
			}

			dest := make([]int, len(randomArr))

			startTime := time.Now()
			segmentMerge(inputA, inputB, dest, threads)
			elapsedTime := time.Since(startTime)

			seqTimeNs := elapsedTime.Nanoseconds()
			mergeDataNs[s][i] = int(seqTimeNs)
		}
	}

	return mergeDataNs
}

func sortBenchmark(sizes []int, iter int) {
	fmt.Printf("Basic merge sort vs fully parallel merge sort at 1, 3 and 6 threads. (Runtimes in ns)\n")
	fmt.Printf("%d Iterations\t\t   Basic\t\t   p = 1\t\t   p = 3\t\t   p = 6\n", iter)

	tabulateResults(sizes, [][][]int{basicSortBenchmark(sizes, iter), parallelSortBenchmark(sizes, iter, 1), parallelSortBenchmark(sizes, iter, 3), parallelSortBenchmark(sizes, iter, 6)})
}

func mergeBenchmark(sizes []int, iter int) {
	fmt.Printf("Sequential merge vs segment merge at p = 1, 3 and 6. (Runtimes in ns)\n")
	fmt.Printf("%d Iterations\t      Sequential\t\t   p = 1\t\t   p = 3\t\t   p = 6\n", iter)

	tabulateResults(sizes, [][][]int{sequentialMergeBenchmark(sizes, iter), segmentMergeBenchmark(sizes, iter, 1), segmentMergeBenchmark(sizes, iter, 3), segmentMergeBenchmark(sizes, iter, 6)})
}

func main() {
	mergeBenchmark([]int{1000, 2000, 4000, 8000, 16000, 32000, 64000, 128000, 256000, 512000, 1024000}, 100)
	fmt.Println()
	sortBenchmark([]int{1000, 2000, 4000, 8000, 16000, 32000, 64000, 128000, 256000}, 100)

}
