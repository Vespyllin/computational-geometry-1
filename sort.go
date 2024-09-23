package main

import "sync"

func segmentMergeSort(I, S []int, p int) {
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
		segmentMergeSort(I[:mid], S[:mid], p)
	}()

	// Sort right half
	go func() {
		defer wg.Done()
		segmentMergeSort(I[mid:], S[mid:], p)
	}()

	wg.Wait()

	// Store the results in the input and copy them back to the scratch space
	copy(I, segmentMerge(S[:mid], S[mid:], p))
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
