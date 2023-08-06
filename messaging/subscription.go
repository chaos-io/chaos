package messaging

const PullMaxWait = 128

type Subscription struct {
	Subject string
	Durable string
	Queue   string

	Pull    bool
	AckWait int
	AutoAck bool
}
