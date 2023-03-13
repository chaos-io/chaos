package trylock

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const mutexLocked = 1 << iota

type Mutex struct {
	in sync.Mutex
}

func (m *Mutex) Lock() {
	m.in.Lock()
}

func (m *Mutex) Unlock() {
	m.in.Unlock()
}

func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.in)), 0, mutexLocked)
}

type RWMutex struct {
	in sync.RWMutex
}

func (rm *RWMutex) RLock() {
	rm.in.RLock()
}

func (rm *RWMutex) RUnlock() {
	rm.in.RUnlock()
}

func (rm *RWMutex) RWLock() {
	rm.in.Lock()
}

func (rm *RWMutex) RWUnlock() {
	rm.in.Unlock()
}

func (rm *RWMutex) TryRWLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&rm.in)), 0, mutexLocked)
}
