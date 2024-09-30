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

func segmentMerge(A []int, B []int, P int) []int {
	var wg sync.WaitGroup

	P = min(P, len(B)) // Limit # of processes so that each process has at least 1 element to sort
	R := make([]int, P+1)
	R[0] = 0

	// Compute ranks
	wg.Add(P)
	for i := 1; i <= P; i++ {
		go func(i int) {
			defer wg.Done()

			idx := min(len(B)-1, (i*len(B))/P) // Cap the target to not overflow B
			R[i] = rank(B[idx], A)
		}(i)
	}

	wg.Wait()

	M := make([]int, len(A)+len(B))
	// Merge each segment
	wg.Add(P)
	for i := 1; i <= P; i++ {
		go func(i int) {
			defer wg.Done()

			bLow := ((i - 1) * (len(B) + 1)) / P
			bHigh := min(((i)*(len(B)+1))/P, len(B)) // Cap the slice to not overflow B

			bSub := B[bLow:bHigh]
			aSub := A[R[i-1]:R[i]]

			// fmt.Printf("A[%v:%v]: %v\n", ranks[i-1], ranks[i], aSub)
			// fmt.Printf("B[%v:%v]: %v\n", bLow, bHigh, bSub)

			sequentialMerge(aSub, bSub, M[R[i-1]+bLow:R[i]+bHigh])
		}(i)
	}

	wg.Wait()

	// Append final elements in A in case of a distribution where the final element in B is smaller than than in A
	for i := R[len(R)-1]; i < len(A); i++ {
		M[len(B)+i] = A[i]
	}

	return M
}
