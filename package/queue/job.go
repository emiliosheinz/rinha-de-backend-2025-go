package queue

type Job interface {
    Execute() error
}

