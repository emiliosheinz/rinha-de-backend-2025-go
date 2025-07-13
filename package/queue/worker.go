package queue

import (
	"log"
	"sync"
)

type Worker struct {
	id   int
	jobs <-chan JobRunner
	wg   *sync.WaitGroup
}

func NewWorker(id int, jobs <-chan JobRunner, wg *sync.WaitGroup) *Worker {
	return &Worker{id: id, jobs: jobs, wg: wg}
}

func (w *Worker) Start() {
	go func() {
		defer w.wg.Done()
		for runner := range w.jobs {
			if err := runner.Job.Execute(runner.Data); err != nil {
				log.Printf("worker %d: job failed: %v", w.id, err)
			}
		}
	}()
}

func StartWorkerPool(n int, queue Queue) []*Worker {
	var wg sync.WaitGroup
	workers := make([]*Worker, n)
	jobs := queue.GetJobs()

	for i := range n {
		wg.Add(1)
		worker := NewWorker(i, jobs, &wg)
		workers[i] = worker
		worker.Start()
	}

	go func() {
		wg.Wait()
	}()

	return workers
}
