package computerclub

import (
	"errors"
	"fmt"
	"sort"
	"testTaskYADRO/internal/parser"
	"testTaskYADRO/internal/utils"
	"time"
)

type computerClub struct {
	numTables int
	startTime time.Time
	endTime   time.Time
	hourCost  int
	queue     []string
	tables    []table
}

func NewComputerClub(numTables int, startTime time.Time, endTime time.Time, hourCost int) *computerClub {
	cc := computerClub{
		numTables: numTables,
		startTime: startTime,
		endTime:   endTime,
		hourCost:  hourCost,
	}
	cc.tables = make([]table, numTables)
	return &cc
}

func (cc *computerClub) Simulate(events []parser.Event) error {
	fmt.Println(cc.startTime.Format("15:04"))

	for _, event := range events {
		event.Print()
	outer:
		switch event.Id {
		case 1:
			if cc.containClient(event.ClientName) {
				cc.makeEvent13(event.Time, "YouShallNotPass")
				break
			}

			if !(event.Time.After(cc.startTime.Add(-time.Minute)) && event.Time.Before(cc.endTime.Add(time.Minute))) {
				cc.makeEvent13(event.Time, "NotOpenYet")
				break
			}

			cc.queue = append(cc.queue, event.ClientName)

		case 2:
			if !cc.containClient(event.ClientName) {
				cc.makeEvent13(event.Time, "ClientUnknown")
				break
			}

			if !cc.tables[event.TableNum-1].isFree() {
				cc.makeEvent13(event.Time, "PlaceIsBusy")
				break
			}

			if cc.getClientTable(event.ClientName) == -1 {
				cc.queue = utils.RemoveMatching(cc.queue, event.ClientName)
			} else {
				cc.tables[cc.getClientTable(event.ClientName)].leave(event.Time, cc.hourCost)
			}
			cc.tables[event.TableNum-1].sit(event.Time, event.ClientName)

		case 3:
			for _, t := range cc.tables {
				if t.isFree() {
					cc.makeEvent13(event.Time, "ICanWaitNoLonger")
					break outer
				}
			}

			if len(cc.queue) > cc.numTables {
				cc.queue = utils.RemoveMatching(cc.queue, event.ClientName)
				cc.makeEvent11(event.Time, event.ClientName)
			}

		case 4:
			if !cc.containClient(event.ClientName) {
				cc.makeEvent13(event.Time, "ClientUnknown")
				break
			}

			tableInx := cc.getClientTable(event.ClientName)

			if tableInx == -1 {
				cc.queue = utils.RemoveMatching(cc.queue, event.ClientName)
				break
			}

			cc.tables[tableInx].leave(event.Time, cc.hourCost)
			if len(cc.queue) != 0 {
				cc.tables[tableInx].sit(event.Time, cc.queue[0])
				cc.makeEvent12(event.Time, cc.queue[0], tableInx+1)
				cc.queue = cc.queue[1:]
			}

		default:
			return errors.New("Event id unknown")
		}
	}

	var remaining []string

	for inx, table := range cc.tables {
		if !table.isFree() {
			remaining = append(remaining, table.client)
			cc.tables[inx].leave(cc.endTime, cc.hourCost)
		}
	}
	for _, client := range cc.queue {
		remaining = append(remaining, client)
		cc.queue = cc.queue[1:]
	}

	sort.Strings(remaining)
	for _, client := range remaining {
		cc.makeEvent11(cc.endTime, client)
	}

	fmt.Println(cc.endTime.Format("15:04"))
	cc.printTableStat()
	return nil
}

func (cc *computerClub) printTableStat() {
	for i, table := range cc.tables {
		hours := int(table.workTime.Hours())
		minutes := int(table.workTime.Minutes()) % 60
		fmt.Printf("%d %d %02d:%02d\n", i+1, table.income, hours, minutes)
	}
}

func (cc *computerClub) getClientTable(client string) int {
	for inx, table := range cc.tables {
		if table.client == client {
			return inx
		}
	}
	return -1
}

func (cc *computerClub) containClient(client string) bool {
	for _, table := range cc.tables {
		if table.client == client {
			return true
		}
	}
	for _, c := range cc.queue {
		if c == client {
			return true
		}
	}
	return false
}

func (cc *computerClub) makeEvent11(time time.Time, clientName string) {
	fmt.Printf("%s %d %s\n", time.Format("15:04"), 11, clientName)
}

func (cc *computerClub) makeEvent12(time time.Time, clientName string, tableNumber int) {
	fmt.Printf("%s %d %s %d\n", time.Format("15:04"), 12, clientName, tableNumber)
}

func (cc *computerClub) makeEvent13(time time.Time, message string) {
	fmt.Printf("%s %d %s\n", time.Format("15:04"), 13, message)
}
