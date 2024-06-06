package tools

import "time"

func GetNow() string {
	now := time.Now().Format("2006-01-02 15:04:05")
	return now
}

// Timer is to call a function periodically.
type Timer struct {
	Function func()
	Duration time.Duration
	Times    int
}

// Start a timer.
func (t *Timer) Start() {
	ticker := time.NewTicker(t.Duration)
	if t.Times > 0 {
		for i := 0; i < t.Times; i++ {
			<-ticker.C
			t.Function()
		}
	} else {
		for range ticker.C {
			t.Function()
		}
	}
}
