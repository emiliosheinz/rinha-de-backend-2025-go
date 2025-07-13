package queue

import (
	"encoding/json"
	"fmt"
)

type InMemoryQueue struct {
	jobs chan JobRunner
}

func NewInMemoryQueue(size int) *InMemoryQueue {
	return &InMemoryQueue{
		jobs: make(chan JobRunner, size),
	}
}

func (q *InMemoryQueue) Enqueue(jobType string, payload any) error {
	job, ok := GetJob(jobType)
	if !ok {
		return fmt.Errorf("unknown job type: %s", jobType)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	q.jobs <- JobRunner{Job: job, Data: data}
	return nil
}

func (q *InMemoryQueue) GetJobs() <-chan JobRunner {
	return q.jobs
}

func (q *InMemoryQueue) Close() {
	close(q.jobs)
}
