//go:build !solution

package once

import (
	"errors"
	"fmt"
)

type PanicInFErr struct {
	panicMsg any
}

func (p *PanicInFErr) Error() string {
	return fmt.Sprint(p.panicMsg)
}

// Once describes an object that will perform exactly one action.
type Once struct {
	blockChan chan struct{}
	done      chan struct{}
}

// New creates Once.
func New() *Once {
	return &Once{
		blockChan: make(chan struct{}),
		done:      make(chan struct{}),
	}
}

// Do calls the function f if and only if Do is being called for the
// first time for this instance of Once. In other words, given
//
//	once := New()
//
// if once.Do(f) is called multiple times, only the first call will invoke f,
// even if f has a different value in each invocation. A new instance of
// Once is required for each function to execute.
//
// Do is intended for initialization that must be run exactly once.
//
// Because no call to Do returns until the one call to f returns, if f causes
// Do to be called, it will deadlock.
//
// If f panics, Do considers it to have returned; future calls of Do return
// without calling f.
func (o *Once) Do(f func()) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				var p *PanicInFErr
				if errors.As(err, &p) {
					// паника произошла при вызове f
					panic(p.panicMsg)
				}
				// паника произошла при повторном закрытии канала(повторный вызов Do), ничего не делаем
			}
			// если конкуррентно вызвали Do ждём завершения первого f
			<-o.done
		}
	}()
	close(o.blockChan)

	defer close(o.done)

	defer func() {
		if r := recover(); r != nil {
			p := &PanicInFErr{
				panicMsg: r,
			}
			panic(fmt.Errorf("panic in f(): %w", p))
		}
	}()

	f()
}
