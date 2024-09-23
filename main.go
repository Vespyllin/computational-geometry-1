package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

// Binary search
func rank(target int, arr []int) int {
	left := 0
	right := len(arr) - 1

	mid := -1
	for left <= right {
		mid = left + (right-left)/2
		// fmt.Printf("L:%d R:%d M:%d V:%d\n", left, right, mid, arr[mid])

		if arr[mid] < target {
			left = mid + 1
		} else if arr[mid] > target {
			right = mid - 1
		} else {
			break // element found
		}

	}

	// fmt.Printf("Final state\n")
	// fmt.Printf("\tT:%d -- L:%d R:%d M:%d V:%d\n", target, left, right, mid, arr[mid])

	// Return index at which to insert at, shifting all elements to the right
	if target <= arr[mid] {
		return mid
	} else {
		return mid + 1
	}
}

func sequentialMerge(a []int, b []int, m []int) {
	i := 0
	j := 0

	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			m[i+j] = a[i]
			i++
		} else {
			m[i+j] = b[j]
			j++
		}
	}

	for i < len(a) {
		m[i+j] = a[i]
		i++
	}

	for j < len(b) {
		m[i+j] = b[j]
		j++
	}
}
func segmentMerge(a []int, b []int, processes int) []int {
	var waitGroup sync.WaitGroup

	processes = min(processes, len(b)) // Limit # of processes so that each process has at least 1 element to sort
	ranks := make([]int, processes+1)
	ranks[0] = 0

	for i := 1; i <= processes; i++ {
		waitGroup.Add(1)
		go func(i int) {
			defer waitGroup.Done()

			idx := min(len(b)-1, (i*len(b))/processes)
			ranks[i] = rank(b[idx], a)
		}(i)
	}
	waitGroup.Wait()

	res := make([]int, len(a)+len(b))
	for i := 1; i <= processes; i++ {
		waitGroup.Add(1)
		go func(i int) {
			defer waitGroup.Done()

			bLow := ((i - 1) * (len(b) + 1)) / processes
			bHigh := min(bLow+(len(b)+1)/processes, len(b))
			bSub := b[bLow:bHigh]

			aSub := a[ranks[i-1]:ranks[i]]

			sequentialMerge(aSub, bSub, res[ranks[i-1]+bLow:ranks[i]+bHigh])
		}(i)
	}

	waitGroup.Wait()

	// Append final elements in A in case of an uneven distribution
	for i := ranks[len(ranks)-1]; i < len(a); i++ {
		res[len(b)+i] = a[i]
	}

	return res
}

func segmentMergeSort(I, S []int) {
	// Base case: if the slice has 0 or 1 element, it's already sorted
	n := len(I)
	if n <= 1 {
		if n == 1 {
			S[0] = I[0]
		}
		return
	}

	mid := n / 2

	var wg sync.WaitGroup
	wg.Add(2)

	// Sort left half
	go func() {
		defer wg.Done()
		basicMergeSort(I[:mid], S[:mid])
	}()

	// Sort right half
	go func() {
		defer wg.Done()
		basicMergeSort(I[mid:], S[mid:])
	}()

	wg.Wait()

	// Store the results in the input and copy them back to the scratch space
	x := segmentMerge(S[:mid], S[mid:], len(S[mid:]))
	copy(I, x)
	copy(S, I)
}

func basicMergeSort(I, S []int) {
	// Base case: if the slice has 0 or 1 element, it's already sorted
	n := len(I)
	if n <= 1 {
		if n == 1 {
			S[0] = I[0]
		}
		return
	}

	mid := n / 2

	var wg sync.WaitGroup
	wg.Add(2)

	// Sort left half
	go func() {
		defer wg.Done()
		basicMergeSort(I[:mid], S[:mid])
	}()

	// Sort right half
	go func() {
		defer wg.Done()
		basicMergeSort(I[mid:], S[mid:])
	}()

	wg.Wait()

	sequentialMerge(S[:mid], S[mid:], I)

	// Copy the sorted result back to the scratch space
	copy(S, I)
}

func genRandomArray(size int) []int {
	n := make([]int, size)

	for i := 0; i < size; i++ {
		n[i] = rand.Int() % 100
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
		fmt.Println("Invalid arguments. Please provide an integer for <inputSize>.")
		return
	}

	threadCount, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid arguments. Please provide an integer for <threadCount>.")
		return
	}

	randomInput := genRandomArray(inputSize)

	if testParam == "seq_merge" {
		fmt.Println("Generating inputs.")

		m := make([]int, inputSize/2)
		a := randomInput[:inputSize/2]
		segmentMergeSort(a, m)

		b := randomInput[inputSize/2:]
		segmentMergeSort(b, m)

		fmt.Println("Starting test.")

		startTime := time.Now()
		sequentialMerge(a, b, m)
		elapsed := time.Since(startTime)

		fmt.Printf("Time elapsed: %s\n", elapsed)
	} else if testParam == "seg_merge" {
		fmt.Println("Generating inputs.")

		// a := []int{4, 5, 6, 8, 10}
		// b := []int{1, 2, 3, 7, 9}
		// a := []int{4, 5, 14, 27, 29, 30, 55, 70}
		// b := []int{10, 16, 17, 25, 28, 40, 59, 80}

		// m1 := make([]int, inputSize/2)
		// m2 := make([]int, inputSize/2)

		a := randomInput[:inputSize/2]
		b := randomInput[inputSize/2:]

		sort.Ints(a)
		sort.Ints(b)

		// segmentMergeSort(a, m1)
		// segmentMergeSort(b, m2)

		fmt.Println("Starting test.")

		startTime := time.Now()
		segmentMerge(a, b, threadCount)
		elapsed := time.Since(startTime)

		fmt.Printf("Time elapsed: %s\n", elapsed)

	} else if testParam == "seq_sort" {
		// Implementation for seq_sort

	} else if testParam == "par_sort" {
		// Implementation for par_sort

	} else if testParam == "seg_sort" {
		// Implementation for seg_sort

	} else {
		fmt.Println("Invalid arguments. Please provide one of seq_merge | seg_merge | seq_sort | par_sort | seg_sort for <testParam>.")
		return
	}

}
