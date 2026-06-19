//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

var ErrStopped = errors.New("limiter stopped")

// request — это запрос от одной горутины на вход в Acquire.
// Ответ управляющей горутины приходит обратно через канал resp.
type request struct {
	ctx  context.Context
	resp chan error
}

// Limiter is precise rate limiter with context support.
type Limiter struct {
	requests chan request
	stop     chan struct{}
	stopped  chan struct{}
}

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxCount events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	l := &Limiter{
		requests: make(chan request),
		stop:     make(chan struct{}),
		stopped:  make(chan struct{}),
	}
	go l.run(maxCount, interval)
	return l
}

// run — единственная управляющая горутина.
//
// Идея: храним кольцевой буфер меток времени последних maxCount успешных
// выдач. Когда приходит новый запрос, смотрим на самую раннюю метку в буфере:
//   - если она старше, чем interval, слот свободен — отвечаем сразу.
//   - иначе надо подождать ровно столько, чтобы та метка «протухла».
//
// Пока ждём — слушаем ctx.Done() запроса и stop.
func (l *Limiter) run(maxCount int, interval time.Duration) {
	defer close(l.stopped)

	timestamps := make([]time.Time, maxCount)
	pos := 0

	for {
		select {
		case <-l.stop:
			return

		case req := <-l.requests:
			if err := req.ctx.Err(); err != nil {
				req.resp <- err
				continue
			}

			oldest := timestamps[pos]
			waitUntil := oldest.Add(interval)
			now := time.Now()

			if waitUntil.Before(now) {
				// Слот уже освободился — выдаём разрешение немедленно.
				timestamps[pos] = now
				pos = (pos + 1) % maxCount
				req.resp <- nil
				continue
			}

			// Нужно подождать delay = waitUntil - now.
			delay := waitUntil.Sub(now)
			timer := time.NewTimer(delay)

			select {
			case <-timer.C:
				// Время вышло — выдаём разрешение.
				timestamps[pos] = time.Now()
				pos = (pos + 1) % maxCount
				req.resp <- nil

			case <-req.ctx.Done():
				// Контекст отменён во время ожидания.
				timer.Stop()
				req.resp <- req.ctx.Err()

			case <-l.stop:
				// Лимитер остановлен во время ожидания.
				timer.Stop()
				req.resp <- ErrStopped
				return
			}
		}
	}
}

// Acquire блокируется до тех пор, пока не будет получено разрешение,
// либо пока не отменится ctx, либо пока не будет вызван Stop.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-l.stopped:
		return ErrStopped
	default:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	resp := make(chan error, 1)
	req := request{ctx: ctx, resp: resp}

	// Отправляем запрос управляющей горутине.
	select {
	case l.requests <- req:
		// Запрос принят; ждём ответа.
		return <-resp
	case <-ctx.Done():
		return ctx.Err()
	case <-l.stopped:
		return ErrStopped
	}
}

func (l *Limiter) Stop() {
	select {
	case <-l.stopped:
	case l.stop <- struct{}{}:
		<-l.stopped
	}
}
