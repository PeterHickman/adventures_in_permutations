package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/etnz/permute"
)

var upto int = 0
var first_element int = 0

func init() {
	n := flag.Int("number", 0, "Number of elements to permute")
	f := flag.Int("first", 0, "The value of the first element")
	flag.Parse()

	if *n == 0 {
		fmt.Printf("The number of elements to permute >= 1 and <= 20")
		os.Exit(1)
	}

	upto = *n
	first_element = *f
}

func main() {
	x := []string{}
	for i := 0; i < upto; i++ {
		x = append(x, strconv.Itoa(i+first_element))
	}

	for _, p := range permute.Permutations(x) {
		fmt.Println(strings.Join(p, ","))
	}
}
