package queue

type Queue interface {
	Enqueue(jobType string, payload any) error
	GetJobs() <-chan JobRunner
	Close()
}
