package messaging

type Message struct {
	Id      string
	TraceId string
	SpanId  string
	Data    string
}
