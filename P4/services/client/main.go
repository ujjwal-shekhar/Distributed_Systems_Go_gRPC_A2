package main

import "flag"

func main() {
	N := flag.Int("N", 1, "Number of Voters")
	T := flag.Int("T", 1, "Number of Traitors")
	flag.Parse()
	if *N < 4 {
		panic("Number of Voters must be at least 4")
	}
	if *T < 0 {
		panic("Number of Traitors must be at least 1")
	}
	// if *N <= *T*3 {
	// 	panic("Number of Voters must be more than 3 times the number of Traitors")
	// }

	StartSimulation(*N, *T)
}