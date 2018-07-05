package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mndrix/rand"
)

func main() {

	fmt.Println(` 
		The Monty Hall Problem Simulation!
		You can specify the number of simulations and number of workers in command line via:
		:>go-montyhall.exe 100 5   (100 simulations, 5 workers)
		
		`)

	simulations := 1000
	workercount := 5

	commLineArgs := os.Args[1:]
	if len(commLineArgs) == 2 {
		argSim := os.Args[1]
		argWrk := os.Args[2]
		simulations, _ = strconv.Atoi(argSim)
		workercount, _ = strconv.Atoi(argWrk)
	}
	jobs := make(chan int, simulations)
	results := make(chan int, simulations)

	swap := true

	// we have two cases to test for, one where the person swapped doors and again when they didn't
	for c := 1; c <= 2; c++ {
		correct := 0
		now := time.Now()
		fmt.Printf("Running Please Wait.. Using %v Workers\n", workercount)
		for w := 1; w <= workercount; w++ {
			go sim(w, jobs, swap, results)
		}

		for j := 1; j <= simulations; j++ {
			jobs <- j
		}

		close(jobs)

		for r := 1; r <= simulations; r++ {
			correct = correct + <-results
			//fmt.Println(<-results)
		}

		timesince := time.Since(now).Seconds()

		fmt.Printf("Did these people swap? %v. \n", swap)
		fmt.Printf("Ran %v Simulations.  Took %.2v seconds\n", simulations, timesince)
		fmt.Printf("They won %v percentage of the time \n\n", PercentOf(correct, simulations))
		// re-run simulation again where the person didn't swap
		close(results)
		swap = false
		jobs = make(chan int, simulations)
		results = make(chan int, simulations)

	}

}

func sim(id int, jobs <-chan int, s bool, results chan<- int) {
	for range jobs {
		winner := rand.Intn(3)
		doors := []int{0, 0, 0}
		doors[winner] = 1
		selected := rand.Intn(3)

		/*
			at this point doors likely looks something like this:   doors := []int{1,0,0}
			the winner has been selected as door[0]
			maybe the randomly selected door by the player is door[1]
			we are going to then assume that door[2] has been opened and shown to the player as garbage.
			we no longer have to consider that doorin the future.
		*/

		switch swap := s; swap {
		case true:
			// if we swap doors and our initially selected door was a loser than we return a winning 1
			if doors[selected] == 0 {
				results <- 1
				//fmt.Println("winner")
			} else {
				//fmt.Println("loser")
				results <- 0
			}

		case false:
			// if we did not swap and the player did not select the winner then they lose.
			if doors[selected] == 0 {
				results <- 0
				//fmt.Println("loser")
			} else {
				//fmt.Println("winner")
				results <- 1

			}

		}

	}
}

func PercentOf(current int, all int) float64 {
	percent := (float64(current) * float64(100)) / float64(all)
	return percent
}
