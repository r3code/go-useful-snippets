// Way to try to enqueue a Job and report back if queue is full and its rejected to add new Job
func TryEnqueue(job Job, jobChan <-chan Job) bool {
    select {
    case jobChan <- job:
        return true
    default:
        return false
    }
}

// Usage
if !TryEnqueue(job, chan) {
    http.Error(w, "max capacity reached", 503)
    return
}
