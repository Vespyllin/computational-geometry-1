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
		fmt.Printf("%-8d\t\t", sizes[j])
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

func main() {
	iter := 100

	mergeSizes := []int{1000, 2000, 4000, 8000, 16000, 32000, 64000, 128000, 256000, 512000, 1024000, 2048000, 4096000, 8192000, 16384000, 32768000, 65536000}
	fmt.Print("Benchmarking sequential merge.")
	sequential := sequentialMergeBenchmark(mergeSizes, iter)
	fmt.Print("\rBenchmarking segment merge (p = 1).")
	segment1 := segmentMergeBenchmark(mergeSizes, iter, 1)
	fmt.Print("\rBenchmarking segment merge (p = 3).")
	segment3 := segmentMergeBenchmark(mergeSizes, iter, 3)
	fmt.Print("\rBenchmarking segment merge (p = 6).")
	segment6 := segmentMergeBenchmark(mergeSizes, iter, 6)
	fmt.Print("\rFinished merge benchmark.")

	fmt.Printf("Sequential merge vs segment merge at p = 1, 3 and 6. %d Iterations (Runtimes in ns)\n", iter)
	fmt.Printf("Size\t\t\t      Sequential\t\t   p = 1\t\t   p = 3\t\t   p = 6\n")

	tabulateResults(mergeSizes, [][][]int{sequential, segment1, segment3, segment6})

	fmt.Println()

	sortSizes := []int{1000, 2000, 4000, 8000, 16000, 32000, 64000, 128000, 256000, 512000, 1024000, 2048000, 4096000}
	fmt.Print("Starting sort benchmark.")
	fmt.Print("\rBenchmarking basic sort.")
	basic := basicSortBenchmark(sortSizes, iter)
	fmt.Print("\rBenchmarking parallel sort (p = 1).")
	parallel1 := parallelSortBenchmark(sortSizes, iter, 1)
	fmt.Print("\rBenchmarking parallel sort (p = 3).")
	parallel3 := parallelSortBenchmark(sortSizes, iter, 3)
	fmt.Print("\rBenchmarking parallel sort (p = 6).")
	parallel6 := parallelSortBenchmark(sortSizes, iter, 6)
	fmt.Print("\rFinished sort benchmark.")

	fmt.Printf("Basic merge sort vs fully parallel merge sort at 1, 3 and 6 threads. %d Iterations (Runtimes in ns)\n", iter)
	fmt.Printf("Size\t\t\t\t   Basic\t\t   p = 1\t\t   p = 3\t\t   p = 6\n")

	tabulateResults(sortSizes, [][][]int{basic, parallel1, parallel3, parallel6})
}
