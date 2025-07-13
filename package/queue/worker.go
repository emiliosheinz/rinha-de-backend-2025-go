package queue

import (
	"fmt"
	"sync"
)

type Worker struct {
	id int
	q  Queue
	wg *sync.WaitGroup
}

func NewWorker(id int, q Queue, wg *sync.WaitGroup) *Worker {
	return &Worker{id: id, q: q, wg: wg}
}

func (w *Worker) Start() {
	go func() {
		defer w.wg.Done()
		for runner := range w.q.GetJobs() {
			if err := runner.Job.Execute(runner.Data); err != nil {
				err := w.q.enqueueRaw(runner.Job.GetType(), runner.Data)
				if err != nil {
					fmt.Printf("Worker %d failed to re-enqueue job %s: %v\n", w.id, runner.Job.GetType(), err)
				} else {
					fmt.Printf("Worker %d re-enqueued job %s\n", w.id, runner.Job.GetType())
				}
			}
		}
	}()
}

func StartWorkerPool(n int, queue Queue) []*Worker {
	var wg sync.WaitGroup
	workers := make([]*Worker, n)

	for i := range n {
		wg.Add(1)
		worker := NewWorker(i, queue, &wg)
		workers[i] = worker
		worker.Start()
	}

	go func() {
		wg.Wait()
	}()

	return workers
}
