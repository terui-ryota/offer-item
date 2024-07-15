package time

import "time"

type Ticker interface {
	Start(interval Duration)
	Tick() <-chan time.Time
	Stop()
}

type ticker struct {
	t        *time.Ticker
	duration time.Duration
}

func NewTicker() Ticker {
	return &ticker{}
}

func (t *ticker) Start(duration Duration) {
	t.t = time.NewTicker(duration.Duration)
}

func (t *ticker) Tick() <-chan time.Time {
	if t.t == nil {
		return nil
	}
	return t.t.C
}

func (t *ticker) Stop() {
	t.t.Stop()
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
