package queue

type Queue interface {
	Enqueue(jobType string, payload any) error
	enqueueRaw(jobType string, payload []byte) error
	GetJobs() <-chan JobRunner
	Close()
}
