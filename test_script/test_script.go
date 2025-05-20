package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/etnz/permute"
)

const target = 12
const number_set = 6

var combs [][]int

func make_comb(size, set int) {
	comb := make([]int, size)

	for {
		t := 0

		// Count the number of values set
		for _, v := range comb {
			if v == 1 {
				t++
			}
		}

		if t == size {
			break
		} else if t == set {
			combs = append(combs, slices.Clone(comb))
		}

		// Step the next entry
		for i, v := range comb {
			if v == 1 {
				comb[i] = 0
			} else {
				comb[i] = 1
				break
			}
		}
	}
}

func merge(left, right, comb []int) {
	r := make([]int, target)

	li := 0
	ri := 0
	for i, v := range comb {
		if v == 0 {
			r[i] = left[li]
			li++
		} else {
			r[i] = right[ri]
			ri++
		}
	}

	// Do the TSP calculation here!
}

func main() {
	start := time.Now()

	make_comb(target, number_set)

	left := []int{1, 2, 3, 4, 5, 6}
	for _, l := range permute.Permutations(left) {
		right := []int{7, 8, 9, 10, 11, 12}
		for _, r := range permute.Permutations(right) {
			for _, c := range combs {
				merge(l, r, c)
			}
		}
	}

	seconds := time.Since(start).Seconds()
	fmt.Printf("Permuted %d elements in %.8f seconds\n", target, seconds)
}
