package main

import (
	// "fmt"
	"sync"
)

// Binary search
func rank(target int, arr []int) int {
    left := 0
	right := len(arr)-1

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
	if (target <= arr[mid]) {
		return mid
	} else {
		return mid + 1
	}
}

func seq_merge(a []int, b []int, m []int) {
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

func seg_merge(a []int, b []int, p int) []int  {
	var wg sync.WaitGroup
	bStep := (len(b) + p - 1) / p // Rounded up step to account for int division imprecision
	
	
	r := make([]int, p+1)
	for i := 1; i <= p; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			target := b[min((bStep*i)-1, len(b)-1)]
			r[i] = rank(target, a)
			// fmt.Println("Target:", target, "Rank:", r[i])
		}(i)
	}

	wg.Wait()

	// fmt.Println("Ranks:",r)
	
	m := make([]int, len(a) + len(b))
	for i := 1; i <= p; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			bLow := (i-1)*bStep 
			bHigh := min(i*bStep, len(b)) // crop the last slice 
			bSub := b[bLow:bHigh]
			// fmt.Printf("B[%d:%d] = %v\n", bLow, bHigh, bSub)
			
			aLow := r[i-1]
			aHigh := r[i]
			aSub := a[aLow:aHigh]		
			// fmt.Printf("A[%d:%d] = %v\n", aLow, aHigh, aSub)

			seq_merge(bSub, aSub, m[aLow+bLow:aHigh+bHigh])
		}(i)
	}

	wg.Wait()

	return m
}



// func main() {
// 	a := []int{4,5,14,27,29,30,55,70}
// 	b := []int{10,16,17,25,28,40,59,80}
//
// 	fmt.Println(seg_merge(a,b,1))
// }
