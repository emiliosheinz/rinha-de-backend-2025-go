package queue

type Job interface {
	Execute(data []byte) error
}

type JobRunner struct {
	Job  Job
	Data []byte
}

var jobRegistry = map[string]Job{}

func RegisterJob(typeName string, job Job) {
	jobRegistry[typeName] = job
}

func GetJob(typeName string) (Job, bool) {
	job, ok := jobRegistry[typeName]
	return job, ok
}
