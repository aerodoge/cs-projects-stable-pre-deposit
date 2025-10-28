package trigger

import (
	"time"
)

type TimeTrigger struct {
	targetTime time.Time
}

func NewTimeTrigger(t time.Time) *TimeTrigger {
	return &TimeTrigger{
		targetTime: t,
	}
}

func (b *TimeTrigger) Listen() (chan time.Time, chan error, error) {
	ch := make(chan time.Time)
	dt := time.Until(b.targetTime)

	timer := time.NewTimer(dt)
	go func() {
		data := <-timer.C
		ch <- data
		close(ch)
	}()

	return ch, nil, nil
}
