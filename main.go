package main

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"time"
)

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

func tabulateResults(sizes []int, A [][]int, B [][]int) {
	for i := 0; i < len(sizes); i++ {
		sumA := 0
		sumB := 0
		for j := 0; j < len(A[i]); j++ {
			sumA += A[i][j]
			sumB += B[i][j]
		}
		// Calculate average
		averageA := float64(sumA) / float64(len(A[i]))
		averageB := float64(sumB) / float64(len(B[i]))

		// Print the row and its average
		comparison := '='
		comparisonColor := "\033[0m"

		if averageA > averageB {
			comparison = '>'
			comparisonColor = "\033[32m"
		} else {
			comparison = '<'
			comparisonColor = "\033[31m"
		}

		fmt.Printf("Size: %s%d\t\t%10.2f\t%2c\t%-10.2f\n\033[0m", comparisonColor, sizes[i], averageA, comparison, averageB)
	}
}

func mergeBenchmark(sizes []int, iter int, threads int) {
	var startTime time.Time
	var elapsedTime time.Duration

	// MERGE
	sequentialDataNs := make([][]int, len(sizes))
	for i := range sizes {
		sequentialDataNs[i] = make([]int, iter)
	}
	segmentDataNs := make([][]int, len(sizes))
	for i := range sizes {
		segmentDataNs[i] = make([]int, iter)
	}

	for s, size := range sizes {
		for i := 0; i < iter; i++ {
			var inputA, inputB []int

			randomArr := genRandUniqueArr(size)
			sortedA := randomArr[len(randomArr)/2:]
			sortedB := randomArr[:len(randomArr)/2]
			sort.Ints(sortedA)
			sort.Ints(sortedB)

			// SEGMENT MERGE
			inputA = make([]int, len(sortedA))
			inputB = make([]int, len(sortedB))
			copy(inputA, sortedA)
			copy(inputB, sortedB)

			rB := rank(inputB[len(inputB)-1], inputA)
			rA := rank(inputA[len(inputA)-1], inputB)
			if rB < len(inputA)/4 || rA < len(inputB)/4 {
				fmt.Printf("Bad distribution at iteration: %d of size %d", iter, size)
				i--
				continue
			}

			startTime = time.Now()
			_ = segmentMerge(inputA, inputB, threads)

			elapsedTime = time.Since(startTime)
			segTimeNs := elapsedTime.Nanoseconds()
			segmentDataNs[s][i] = int(segTimeNs)

			// SEQUENTIAL MERGE
			scratch := make([]int, len(randomArr))
			inputA = make([]int, len(sortedA))
			inputB = make([]int, len(sortedB))
			copy(inputA, sortedA)
			copy(inputB, sortedB)

			startTime = time.Now()
			sequentialMerge(inputA, inputB, scratch)
			elapsedTime = time.Since(startTime)

			seqTimeNs := elapsedTime.Nanoseconds()
			sequentialDataNs[s][i] = int(seqTimeNs)
		}
	}

	fmt.Printf("Sequential and Segment merge algorithm average runtimes in ns.\n")
	fmt.Printf("%d Iterations\t\tSequential\t\tSegment\n", iter)
	tabulateResults(sizes, sequentialDataNs, segmentDataNs)
}

func sortBenchmark(sizes []int, iter int, threads int) {
	var startTime time.Time
	var elapsedTime time.Duration

	basicSortDataNs := make([][]int, len(sizes))
	for i := range sizes {
		basicSortDataNs[i] = make([]int, iter)
	}
	segmentSortDataNs := make([][]int, len(sizes))
	for i := range sizes {
		segmentSortDataNs[i] = make([]int, iter)
	}

	// var inputA, inputB []int
	for s, size := range sizes {
		input := make([]int, size)
		var scratch []int
		for i := 0; i < iter; i++ {
			randomArr := genRandUniqueArr(size)

			// SEGMENT MERGE
			copy(input, randomArr)
			scratch = make([]int, len(randomArr))

			startTime = time.Now()
			segmentMergeSort(input, scratch, threads)
			elapsedTime = time.Since(startTime)

			if !sort.IntsAreSorted(input) {
				fmt.Printf("ERROR IN SEGMENT MERGE SORT\n\tSIZE: %d - ITERATION: %d - THREADS: %d\n", size, i, threads)
				break
			}
			segTimeNs := elapsedTime.Nanoseconds()
			segmentSortDataNs[s][i] = int(segTimeNs)

			// SEQUENTIAL MERGE
			copy(input, randomArr)
			scratch = make([]int, len(randomArr))

			startTime = time.Now()
			basicMergeSort(input, scratch)
			elapsedTime = time.Since(startTime)

			if !sort.IntsAreSorted(input) {
				fmt.Printf("ERROR IN BASIC MERGE SORT\n\tSIZE: %d - ITERATION: %d - THREADS: %d\n", size, i, threads)
				break
			}

			basicTimeNs := elapsedTime.Nanoseconds()
			basicSortDataNs[s][i] = int(basicTimeNs)
		}
	}

	fmt.Printf("Basic and Segment merge sort algorithm average runtimes in ns.\n")
	fmt.Printf("%d Iterations\t\tBasic\t\t\tSegment\n", iter)
	tabulateResults(sizes, basicSortDataNs, segmentSortDataNs)
}

func main() {
	// mergeBenchmark([]int{1000, 2000, 4000, 8000, 16000, 32000, 64000, 128000, 256000, 512000, 1024000}, 1000, 6)
	sortBenchmark([]int{1000, 2000, 4000, 8000, 16000, 32000, 64000, 128000, 256000}, 100, 6)
}
