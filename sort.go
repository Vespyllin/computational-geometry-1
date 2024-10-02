package main

import (
	"sync"
)

func segmentMergeSort(I, S []int, p int) ([]int, []int) {
	n := len(I)
	if n <= 1 {
		if n == 1 {
			copy(S, I)
		}
		return S, I
	}

	mid := n / 2

	Il, Ir := I[:mid], I[mid:]
	Sl, Sr := S[:mid], S[mid:]

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, Sl = segmentMergeSort(Il, Sl, p)
	}()

	go func() {
		defer wg.Done()
		_, Sr = segmentMergeSort(Ir, Sr, p)
	}()

	wg.Wait()

	res := make([]int, len(S))
	segmentMerge(Sr, Sl, res, p)
	return S, res
}

func basicMergeSort(I, S []int) ([]int, []int) {
	n := len(I)
	if n <= 1 {
		if n == 1 {
			copy(S, I)
		}
		return S, I
	}

	mid := n / 2

	Il, Ir := I[:mid], I[mid:]
	Sl, Sr := S[:mid], S[mid:]

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, Sl = basicMergeSort(Il, Sl)
	}()

	go func() {
		defer wg.Done()
		_, Sr = basicMergeSort(Ir, Sr)
	}()

	wg.Wait()

	res := make([]int, len(S))
	sequentialMerge(Sr, Sl, res)
	return S, res
}
