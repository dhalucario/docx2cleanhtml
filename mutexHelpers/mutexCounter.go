package mutexHelpers

import "sync"

type MutexCounter struct {
	val int
	lock sync.Mutex
}

func NewMutexCounter(n int) MutexCounter {
	return MutexCounter{
		val: n,
		lock: sync.Mutex{},
	}
}

func (mc* MutexCounter) Increment(n int) int {
	mc.lock.Lock()
	defer mc.lock.Unlock()
	mc.val = mc.val + n
	return mc.val
}

func (mc* MutexCounter) Decrement(n int) int {
	mc.lock.Lock()
	defer mc.lock.Unlock()
	mc.val = mc.val - n
	return mc.val
}

func (mc* MutexCounter) LSet(n int) int {
	mc.lock.Lock()
	defer mc.lock.Unlock()
	mc.val = n
	return mc.val
}

func (mc* MutexCounter) LGet() int {
	mc.lock.Lock()
	defer mc.lock.Unlock()
	return mc.val
}

func (mc* MutexCounter) Set(n int) int {
	mc.val = n
	return mc.val
}

func (mc* MutexCounter) Get() int {

	return mc.val
}

func (mc* MutexCounter) Lock() {
	mc.lock.Lock()
}

func (mc* MutexCounter) Unlock() {
	mc.lock.Unlock()
}