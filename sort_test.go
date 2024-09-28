package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// generateUniqueRandomArray generates an array of unique random integers of given size.
func generateUniqueRandomArray(size int) []int {
	rand.Seed(time.Now().UnixNano()) // Seed for randomness using the current time
	m := make(map[int]struct{})      // A map to track unique values (using empty struct to save memory)
	result := make([]int, 0, size)   // Slice to store the result

	// Continue generating random numbers until we have 'size' unique numbers
	for len(result) < size {
		n := rand.Intn(size * 2)        // Generate numbers within a range slightly larger than 'size'
		if _, exists := m[n]; !exists { // Check if the number is already in the map
			m[n] = struct{}{}          // Mark this number as used
			result = append(result, n) // Add the unique number to the result array
		}
	}
	return result
}

// BenchmarkSorts benchmarks the performance of SegmentMergeSort and BasicMergeSort
// across different input sizes and thread counts.
func BenchmarkSorts(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000} // Array sizes to benchmark
	threadCounts := []int{1, 2, 4, 8}        // Different thread counts for SegmentMergeSort

	// Loop through each size and thread count for benchmarking
	for _, size := range sizes {
		input := generateUniqueRandomArray(size) // Generate the input array for this size

		// Benchmark SegmentMergeSort with different thread counts
		for _, threads := range threadCounts {
			b.Run(fmt.Sprintf("SegmentMergeSort-Size%d-Threads%d", size, threads), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					b.StopTimer() // Stop the timer to avoid measuring setup time

					// Prepare a copy of the input array for each run
					testInput := make([]int, len(input))
					copy(testInput, input)
					scratch := make([]int, len(testInput)) // Scratch array for sorting

					b.StartTimer()                                // Start timing the actual sorting
					segmentMergeSort(testInput, scratch, threads) // Perform the sort
				}
			})
		}

		// Benchmark BasicMergeSort (single-threaded)
		b.Run(fmt.Sprintf("BasicMergeSort-Size%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer() // Stop the timer during setup

				// Prepare a copy of the input array for each run
				testInput := make([]int, len(input))
				copy(testInput, input)
				scratch := make([]int, len(testInput)) // Scratch array for sorting

				b.StartTimer()                     // Start timing the actual sorting
				basicMergeSort(testInput, scratch) // Perform the basic merge sort
			}
		})
		fmt.Println()
	}
}

// TestSortCorrectness ensures that both SegmentMergeSort and BasicMergeSort
// are sorting the arrays correctly.
func TestSortCorrectness(t *testing.T) {
	sizes := []int{100, 1000, 10000}  // Array sizes to test
	threadCounts := []int{1, 2, 4, 8} // Different thread counts for SegmentMergeSort

	// Loop through each size and thread count for testing
	for _, size := range sizes {
		input := generateUniqueRandomArray(size) // Generate the input array for this size

		// Test SegmentMergeSort with different thread counts
		for _, threads := range threadCounts {
			t.Run(fmt.Sprintf("SegmentMergeSort-Size%d-Threads%d", size, threads), func(t *testing.T) {
				// Prepare a copy of the input array for each test
				testInput := make([]int, len(input))
				copy(testInput, input)
				scratch := make([]int, len(testInput)) // Scratch array for sorting

				segmentMergeSort(testInput, scratch, threads) // Perform the sort

				// Verify that the array is correctly sorted
				if !isSorted(testInput) {
					t.Errorf("SegmentMergeSort failed to sort the array correctly")
				}
			})
		}

		// Test BasicMergeSort (single-threaded)
		t.Run(fmt.Sprintf("BasicMergeSort-Size%d", size), func(t *testing.T) {
			// Prepare a copy of the input array for each test
			testInput := make([]int, len(input))
			copy(testInput, input)
			scratch := make([]int, len(testInput)) // Scratch array for sorting

			basicMergeSort(testInput, scratch) // Perform the basic merge sort

			// Verify that the array is correctly sorted
			if !isSorted(testInput) {
				t.Errorf("BasicMergeSort failed to sort the array correctly")
			}
		})
	}
}

// isSorted checks if the given array is sorted in ascending order
func isSorted(arr []int) bool {
	// Loop through the array and check that each element is less than or equal to the next
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return false
		}
	}
	return true
}

// TestMain sets up the test environment by seeding the random generator.
// The random seed is set to the current time for randomness in tests.
func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator for tests
	m.Run()                          // Run all tests
}
