package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/etnz/permute"
)

var upto int = 0

func init() {
	n := flag.Int("number", 0, "Number of elements to permute")
	flag.Parse()

	if *n == 0 {
		fmt.Printf("The number of elements to permute >= 1 and <= 20")
		os.Exit(1)
	}

	upto = *n
}

func main() {
	x := []int{}
	for i := 0; i < upto; i++ {
		x = append(x, i)
	}

	start := time.Now()

	for _, _ = range permute.Permutations(x) {
	}

	seconds := time.Since(start).Seconds()
	fmt.Printf("Permuted %d elements in %.8f seconds\n", upto, seconds)
}
