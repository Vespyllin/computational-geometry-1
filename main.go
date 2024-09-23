package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"sort"
	"strconv"
	"time"
)

func genRandomArray(size int) []int {
	n := make([]int, size)

	for i := 0; i < size; i++ {
		n[i] = rand.Int() % 300
	}

	return n
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <testParam: {seq_merge | seg_merge | seq_sort | par_sort | seg_sort}> <inputSize> <threadCount>")
		return
	}

	testParam := os.Args[1]

	inputSize, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid arguments. Provide an integer for <inputSize>.")
		return
	}

	threadCount, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid arguments. Provide an integer for <threadCount>.")
		return
	}

	randomInput := genRandomArray(inputSize)

	if testParam == "test" {
		fmt.Println("Generating inputs.")

		fmt.Println("Starting test.")

		startTime := time.Now()
		a := randomInput[:inputSize/2]
		b := randomInput[inputSize/2:]
		sort.Ints(a)
		sort.Ints(b)

		// m1 := make([]int, len(randomInput))
		// segmentMergeSort(randomInput, m1, threadCount)
		res := segmentMerge(a, b, threadCount)
		fmt.Println(sort.IntsAreSorted(res))
		elapsed := time.Since(startTime)

		fmt.Printf("Time elapsed: %s\n", elapsed)
	} else if testParam == "seq_merge" {
		fmt.Println("Generating inputs.")

		m := make([]int, inputSize/2)
		a := randomInput[:inputSize/2]
		segmentMergeSort(a, m, 8)

		b := randomInput[inputSize/2:]
		segmentMergeSort(b, m, 8)

		fmt.Println("Starting test.")

		startTime := time.Now()
		sequentialMerge(a, b, m)
		elapsed := time.Since(startTime)

		fmt.Printf("Time elapsed: %s\n", elapsed)
	} else if testParam == "seg_merge" {

		// Implementation for seq_merge

	} else if testParam == "seq_sort" {
		// Implementation for seq_sort

	} else if testParam == "par_sort" {
		// Implementation for par_sort

	} else if testParam == "seg_sort" {
		// Implementation for seg_sort

	} else {
		fmt.Println("Invalid arguments. Provide one of seq_merge | seg_merge | seq_sort | par_sort | seg_sort for <testParam>.")
		return
	}

}
