package messaging

//go:generate mockgen -destination=mocks/queue.go -package=mocks . Queue
type Queue interface {
	Publisher
	Subscriber

	Shutdown()
}
