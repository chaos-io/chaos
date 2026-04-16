package messaging

import (
	"sync"
	"sync/atomic"
)

type SubMessage struct {
	Message

	mu         sync.RWMutex
	done       atomic.Bool
	ack        func()
	nak        func()
	term       func()
	inProgress func()
}

func (m *SubMessage) SetAck(ack func()) {
	m.setAction(func() { m.ack = ack })
}

func (m *SubMessage) SetNak(nak func()) {
	m.setAction(func() { m.nak = nak })
}

func (m *SubMessage) SetTerm(term func()) {
	m.setAction(func() { m.term = term })
}

func (m *SubMessage) SetInProgress(inProgress func()) {
	m.setAction(func() { m.inProgress = inProgress })
}

func (m *SubMessage) Ack() {
	m.finish(m.getAction(func(message *SubMessage) func() { return message.ack }))
}

func (m *SubMessage) Nak() {
	m.finish(m.getAction(func(message *SubMessage) func() { return message.nak }))
}

func (m *SubMessage) Term() {
	m.finish(m.getAction(func(message *SubMessage) func() { return message.term }))
}

func (m *SubMessage) InProgress() {
	if fn := m.getAction(func(message *SubMessage) func() { return message.inProgress }); fn != nil {
		fn()
	}
}

func (m *SubMessage) setAction(setter func()) {
	if m == nil {
		return
	}
	m.mu.Lock()
	setter()
	m.mu.Unlock()
}

func (m *SubMessage) getAction(selector func(message *SubMessage) func()) func() {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	fn := selector(m)
	m.mu.RUnlock()
	return fn
}

func (m *SubMessage) finish(fn func()) {
	if m == nil || fn == nil {
		return
	}
	if !m.done.CompareAndSwap(false, true) {
		return
	}
	fn()
}
