package queue

import (
	"log"
	"sync"
)

type Worker struct {
	id   int
	jobs <-chan Job
	wg   *sync.WaitGroup
}

func NewWorker(id int, jobs <-chan Job, wg *sync.WaitGroup) *Worker {
	return &Worker{id: id, jobs: jobs, wg: wg}
}

func NewWorkerPool(num int, jobs <-chan Job, wg *sync.WaitGroup) []*Worker {
	workers := make([]*Worker, num)
	for i := range num {
		wg.Add(1)
		worker := NewWorker(i, jobs, wg)
		workers[i] =  worker
		worker.Start()	
	}
	return workers
}

func (w Worker) Start() {
	go func() {
		defer w.wg.Done()
		for job := range w.jobs {
			if err := job.Execute(); err != nil {
				// What if both tries fail, the job will fail, we need to handle that
				log.Printf("worker %d: job error: %v", w.id, err)
			}
		}
	}()
}
