package main

import "fmt"

func worker(jobChan <-chan Job) {
	for job := range jobChan {
		process(job)
	}
}

func process(j Job) {
	fmt.Printf("job: %s", j.Name)
}

type Job struct {
	Name string
}

func main() {
	// make a channel with a capacity of 100. Allow us to throttle the producer (block acdding new jobs until input queue has an empty slot)
	jobChan := make(chan Job, 100)

	// start the worker
	go worker(jobChan)

	for i := 0; i < 200; i++ {
		job := Job{"a"}
		// enqueue a job
		jobChan <- job
		// will block, if there already are 100 jobs in the channel
	}
}
