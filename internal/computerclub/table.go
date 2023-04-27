package computerclub

import "time"

type table struct {
	currentTime time.Time
	workTime    time.Duration
	client      string
	income      int
}

func (t *table) isFree() bool {
	return t.client == ""
}

func (t *table) sit(time time.Time, name string) {
	t.currentTime = time
	t.client = name
}

func (t *table) leave(time time.Time, hourCost int) {
	t.client = ""

	currentWorkTime := time.Sub(t.currentTime)
	t.workTime = t.workTime + currentWorkTime

	if int(currentWorkTime.Minutes())%60 != 0 {
		t.income += hourCost * (int(currentWorkTime.Hours()) + 1)
	} else {
		t.income += hourCost * int(currentWorkTime.Hours())
	}

	t.currentTime = time
}
