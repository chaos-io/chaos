package messaging

import (
	"sync/atomic"
)

type SubMessage struct {
	Message

	done       atomic.Bool
	ack        func()
	nak        func()
	term       func()
	inProgress func()
}

func (m *SubMessage) SetAck(ack func()) {
	if m == nil {
		return
	}
	m.ack = ack
}

func (m *SubMessage) SetNak(nak func()) {
	if m == nil {
		return
	}
	m.nak = nak
}

func (m *SubMessage) SetTerm(term func()) {
	if m == nil {
		return
	}
	m.term = term
}

func (m *SubMessage) SetInProgress(inProgress func()) {
	if m == nil {
		return
	}
	m.inProgress = inProgress
}

func (m *SubMessage) Ack() {
	if m == nil {
		return
	}
	m.finish(m.ack)
}

func (m *SubMessage) Nak() {
	if m == nil {
		return
	}
	m.finish(m.nak)
}

func (m *SubMessage) Term() {
	if m == nil {
		return
	}
	m.finish(m.term)
}

func (m *SubMessage) InProgress() {
	if m == nil || m.inProgress == nil {
		return
	}
	m.inProgress()
}

func (m *SubMessage) Done() bool {
	if m == nil {
		return false
	}
	return m.done.Load()
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
