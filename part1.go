package main

import (
	"fmt"
	"sync"
)

// BasicMergeSort performs a parallel merge sort on the input slice.
// It takes two slices as input: I (the input array) and S (the scratch space),
// and a pointer to a WaitGroup for synchronization.
func BasicMergeSort(I, S []int, wg *sync.WaitGroup) {
    defer wg.Done()

    n := len(I)
    if n <= 1 {
        // Base case: if the slice has 0 or 1 element, it's already sorted
        if n == 1 {
            S[0] = I[0]
        }
        return
    }

    // Split the input
    mid := n / 2
    var wgHalves sync.WaitGroup
    wgHalves.Add(2)

    // Sort the left half
    go BasicMergeSort(I[:mid], S[:mid], &wgHalves)

    // Sort the right half
    go BasicMergeSort(I[mid:], S[mid:], &wgHalves)

    // Wait for both sorting operations to complete
    wgHalves.Wait()

    // Merge the sorted halves into the scratch space
    merge(S[:mid], S[mid:], I)

    // Copy the sorted result back to the scratch space
    copy(S, I)
}

// merge combines two sorted slices (left and right) into a single sorted result.
// This is a sequential merge operation.
func merge(left, right, result []int) {
    i, j, k := 0, 0, 0

    // Compare elements from both slices and put the smaller one into the result
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result[k] = left[i]
            i++
        } else {
            result[k] = right[j]
            j++
        }
        k++
    }

    // If there are remaining elements in left, append them to result
    for i < len(left) {
        result[k] = left[i]
        i++
        k++
    }

    // If there are remaining elements in right, append them to result
    for j < len(right) {
        result[k] = right[j]
        j++
        k++
    }
}

func main() {
	// Example usage of the BasicMergeSort function
	input := []int{64, 34, 25, 12, 22, 11, 90, 1}
	scratch := make([]int, len(input))
	
	fmt.Println("Original array:", input)

	var wg sync.WaitGroup
	wg.Add(1)
	BasicMergeSort(input, scratch, &wg)
	wg.Wait()

	fmt.Println("Sorted array:", input)
}
