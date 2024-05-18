package runner

import "sync"

type Lock interface {
	Lock()
	Unlock()
	ReadLock()
	ReadUnlock()
}

type MutexLock struct {
	lock sync.RWMutex
}

func NewMutexLock() Lock {
	return &MutexLock{
		lock: sync.RWMutex{},
	}
}

func (m *MutexLock) Lock() {
	m.lock.Lock()
}

func (m *MutexLock) Unlock() {
	m.lock.Unlock()
}

func (m *MutexLock) ReadLock() {
	m.lock.RLock()
}

func (m *MutexLock) ReadUnlock() {
	m.lock.RUnlock()
}
