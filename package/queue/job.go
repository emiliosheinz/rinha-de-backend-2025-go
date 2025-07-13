package queue

type Job interface {
	Execute(data []byte) error
	GetType() string
}

type JobRunner struct {
	Job  Job
	Data []byte
}

var jobRegistry = map[string]Job{}

func RegisterJob(job Job) {
	jobRegistry[job.GetType()] = job
}

func GetJob(typeName string) (Job, bool) {
	job, ok := jobRegistry[typeName]
	return job, ok
}
