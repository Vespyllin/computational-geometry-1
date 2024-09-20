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
    seq_merge(S[:mid], S[mid:], I)

    // Copy the sorted result back to the scratch space
    copy(S, I)
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
