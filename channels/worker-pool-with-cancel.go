
func worker(jobChan <-chan Job, cancelChan <-chan struct{}) {
    for {
        select {
        case <-cancelChan:
            return

        case job := <-jobChan:
            process(job)
        }
    }
}

// create cancel channel
cancelChan := make(chan struct{})

// pass the channel to the workers, let them wait on it
for i:=0; i<workerCount; i++ {
    go worker(jobChan, cancelChan)
}

// close the channel to signal the workers
close(cancelChan)
