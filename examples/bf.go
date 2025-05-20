package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	// "runtime/pprof"

	"github.com/PeterHickman/expand_path"
	"github.com/PeterHickman/toolbox"
	"github.com/etnz/permute"
)

// "./bf --root #{root} --right_size 8 --right_start 8 --tsp tsp_16_cities_symetric.txt

var root []int
var root_size int
var right_size int
var right_start int
var tsp string

var combs [][]int

// TSP Data

const unconnected = 9999999

var best_route = 9999999
var symetrical bool = false
var city_names []string
var city_map map[string]int
var distances []int

// TSP functions

func true_or_false(text string) bool {
	if text == "true" || text == "yes" {
		return true
	}

	if text == "false" || text == "no" {
		return false
	}

	log.Fatal("[" + text + "] is neither true nor false")
	return false
}

// https://www.reddit.com/r/golang/comments/818953/whats_the_cleanest_way_to_fill_up_a_slice_with/
func fill[T any](slice []T, val T) {
	for i := range slice {
		slice[i] = val
	}
}

func get_list_of_cities(edge_list [][]string) []string {
	cities := []string{}

	for _, edge := range edge_list {
		if !slices.Contains(cities, edge[0]) {
			cities = append(cities, edge[0])
		}
		if !slices.Contains(cities, edge[1]) {
			cities = append(cities, edge[1])
		}
	}

	slices.Sort(cities)
	return cities
}

func measure_route(route []int) int {
	t := 0

	l := len(city_names)
	for i := 0; i < len(route)-1; i++ {
		t += distances[route[i]*l+route[i+1]]
	}
	t += distances[route[l-1]*l+route[0]]

	return t
}

func read_file(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	symetrical_is_set := false
	edge_list := make([][]string, 0)

	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		parts := strings.Fields(line)

		if len(parts) > 0 {
			switch parts[0] {
			case "symetrical":
				if symetrical_is_set {
					log.Fatal("symetrical is already set")
				} else if len(parts) == 2 {
					symetrical_is_set = true
					symetrical = true_or_false(parts[1])
				} else {
					log.Fatal("symetrical takes only 1 argument")
				}
			case "fullyconnected":
			case "edge":
				if len(parts) == 4 {
					edge_list = append(edge_list, parts[1:])
				} else {
					log.Fatal("edge takes 3 arguments")
				}
			}
		}
	}

	city_names = get_list_of_cities(edge_list)
	city_map = map[string]int{}

	for i, n := range city_names {
		city_map[n] = i
	}

	l := len(city_names)
	distances = make([]int, l*l)
	fill(distances, unconnected)

	for _, edge := range edge_list {
		a := city_map[edge[0]]
		b := city_map[edge[1]]
		d, _ := strconv.Atoi(edge[2])

		distances[a*l+b] = d
		if symetrical {
			distances[b*l+a] = d
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// End of TSP functions

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
		} else if t == set && comb[0] == 0 {
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

func merge(left, right, comb []int, target int) {
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

	l := measure_route(r)
	if l < best_route {
		best_route = l
	}
}

func parse_root(text string) {
	parts := strings.Split(text, ",")
	root_size = len(parts)

	for _, v := range parts {
		val, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("Expecting an integer in the root. got %s", v)
		}
		if val < 0 || val > 255 {
			log.Fatalf("Root values should be 0 <= x <= 255. got %d", val)
		}
		root = append(root, val)
	}

	// TODO: Additionally check for duplicate values
}

func init() {
	r := flag.String("root", "", "The root, left hand side, of the problem space")
	rs := flag.Int("right_size", 0, "The size of the right hand side")
	rf := flag.Int("right_start", 0, "The lowest value for the elements of the right hand side")
	t := flag.String("tsp", "", "The file defining the TSP problem to solve")

	flag.Parse()

	if *r == "" {
		log.Fatal("You must supply a root with --root")
	}

	parse_root(*r)

	if *rs < 1 {
		log.Fatalf("The size of the right hand side needs to be > 1")
	}

	right_size = *rs

	if *rf < 1 {
		log.Fatalf("The lowest element of the right hand side should be positive")
	}

	for _, v := range root {
		if v == *rf {
			log.Fatalf("The lowest element of the right hand side should not be in the root (left hand side)")
		}
	}

	right_start = *rf

	if *t == "" {
		log.Fatal("You must supply a TSP to solve with --tsp")
	}

	e_tsp, _ := expand_path.ExpandPath(*t)

	if !toolbox.FileExists(e_tsp) {
		log.Fatalf("The TSP file [%s] not found", *t)
	}

	tsp = e_tsp
}

func main() {
	// Start profiling
	// f, err := os.Create("myprogram.prof")
	// if err != nil {
	//   fmt.Println(err)
	//   return
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	started := time.Now()

	make_comb(root_size+right_start, right_start)
	read_file(tsp)

	right := []int{}
	for i := 0; i < right_size; i++ {
		right = append(right, right_start+i)
	}

	for _, r := range permute.Permutations(right) {
		for _, c := range combs {
			merge(root, r, c, root_size+right_start)
		}
	}

	elapsed := time.Since(started).Seconds()

	fmt.Printf("Best route was %d\n", best_route)
	fmt.Printf("Ran in %.6f seconds\n", elapsed)
}
