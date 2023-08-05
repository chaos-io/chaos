package messaging

const PullMaxWait = 128

type Subscribe struct {
	Subject string
	Durable string
	Queue   string
	Pull    bool
	AckWait int
}
