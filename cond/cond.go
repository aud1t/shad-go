//go:build !solution

package cond

// A Locker represents an object that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}

// Cond implements a condition variable, a rendezvous point
// for goroutines waiting for or announcing the occurrence
// of an event.
//
// Each Cond has an associated Locker L (often a *sync.Mutex or *sync.RWMutex),
// which must be held when changing the condition and
// when calling the Wait method.
type Cond struct {
	L       Locker
	waiters []chan struct{}
	mu      chan struct{}
}

// New returns a new Cond with Locker l.
func New(l Locker) *Cond {
	c := &Cond{
		L:  l,
		mu: make(chan struct{}, 1),
	}
	c.mu <- struct{}{}
	return c
}

// Wait atomically unlocks c.L and suspends execution
// of the calling goroutine. After later resuming execution,
// Wait locks c.L before returning. Unlike in other systems,
// Wait cannot return unless awoken by Broadcast or Signal.
//
// Because c.L is not locked when Wait first resumes, the caller
// typically cannot assume that the condition is true when
// Wait returns. Instead, the caller should Wait in a loop:
//
//	c.L.Lock()
//	for !condition() {
//	    c.Wait()
//	}
//	... make use of condition ...
//	c.L.Unlock()
func (c *Cond) Wait() {
	<-c.mu
	waitChan := make(chan struct{})
	c.waiters = append(c.waiters, waitChan)
	c.mu <- struct{}{}

	c.L.Unlock()
	<-waitChan
	c.L.Lock()
}

// Signal wakes one goroutine waiting on c, if there is any.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Signal() {
	<-c.mu
	if len(c.waiters) == 0 {
		c.mu <- struct{}{}
		return
	}
	wchan := c.waiters[0]
	c.waiters = c.waiters[1:]
	c.mu <- struct{}{}

	close(wchan)
}

// Broadcast wakes all goroutines waiting on c.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Broadcast() {
	<-c.mu
	waiters := c.waiters
	c.waiters = nil
	c.mu <- struct{}{}

	for _, waiter := range waiters {
		close(waiter)
	}
}
