package queue

type Queue interface {
	Enqueue(job Job)
	GetJobs() <-chan Job
	Close()
}

