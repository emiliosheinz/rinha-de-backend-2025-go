package queue

type InMemoryQueue struct {
	jobs chan Job
}

func NewInMemoryQueue(size int) *InMemoryQueue {
	return &InMemoryQueue{
		jobs: make(chan Job, size),
	}
}

func (q *InMemoryQueue) GetJobs() <-chan Job {
	return q.jobs
}

func (q *InMemoryQueue) Enqueue(job Job) {
	q.jobs <- job
}

func (q *InMemoryQueue) Close() {
	close(q.jobs)
}
