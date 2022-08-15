package timer

type Timer struct {
	ticks       int
	targetTicks int
}

func NewTimer(target int) *Timer {
	return &Timer{
		ticks:       0,
		targetTicks: target,
	}
}

func (t *Timer) Update() {
	t.ticks++
}

func (t *Timer) IsDone() bool {
	return t.ticks >= t.targetTicks
}

func (t *Timer) Reset() {
	t.ticks = 0
}
