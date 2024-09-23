package main

import (
	"sync"
)

func rank(target int, arr []int) int {
	// Binary search
	left := 0
	right := len(arr) - 1
	mid := -1
	for left <= right {
		mid = left + (right-left)/2

		if arr[mid] < target {
			left = mid + 1
		} else if arr[mid] > target {
			right = mid - 1
		} else {
			break // Bounds located
		}

	}

	// Return index of next greatest element
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
			bHigh := min(((i)*(len(b)+1))/processes, len(b))

			bSub := b[bLow:bHigh]
			aSub := a[ranks[i-1]:ranks[i]]

			// fmt.Printf("A[%v:%v]: %v\n", ranks[i-1], ranks[i], aSub)
			// fmt.Printf("B[%v:%v]: %v\n", bLow, bHigh, bSub)

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
