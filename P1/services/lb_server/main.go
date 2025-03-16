package main

import "flag"

func main () {
	POLICY := flag.String("POLICY", "least_loaded", "The policy to be used by the load balancer")
	flag.Parse()

	if *POLICY != "least_loaded" && *POLICY != "round_robin" && *POLICY != "pick_first" {
		panic("Invalid policy. Please use either 'least_loaded' or 'round_robin'")
	}

	StartLBServer(*POLICY)
}