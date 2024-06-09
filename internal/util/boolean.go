package util

import "sync/atomic"

type AtomicBool struct {
	value int32
}

func (b *AtomicBool) Set(val bool) {
	var intVal int32
	if val {
		intVal = 1
	} else {
		intVal = 0
	}
	atomic.StoreInt32(&b.value, intVal)
}

func (b *AtomicBool) Get() bool {
	return atomic.LoadInt32(&b.value) != 0
}
