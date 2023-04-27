package parser

import (
	"fmt"
	"time"
)

type Event struct {
	Time       time.Time
	Id         int
	ClientName string
	TableNum   int
}

func (event *Event) Print() {
	if event.TableNum != 0 {
		fmt.Printf("%s %d %s %d\n", event.Time.Format("15:04"), event.Id, event.ClientName, event.TableNum)
	} else {
		fmt.Printf("%s %d %s\n", event.Time.Format("15:04"), event.Id, event.ClientName)
	}
}
