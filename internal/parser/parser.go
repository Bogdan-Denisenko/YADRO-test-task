package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func parseTime(t string) (time.Time, error) {
	if len(t) != 5 {
		return time.Time{}, errors.New("Time format not XX:XX")
	}
	parsedTime, err := time.Parse("15:04", t)
	return parsedTime, err
}

func checkName(name string) error {
	const pattern = "^[a-z0-9_-]+$"
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	if !re.MatchString(name) {
		return fmt.Errorf("invalid name format: %s", name)
	}

	return nil
}

func ParseFile(filePath string) (int, time.Time, time.Time, int, []Event, error) {
	var numTables int
	var startTime time.Time
	var endTime time.Time
	var hourCost int
	var events []Event

	file, err := os.Open(filePath)
	if err != nil {
		return 0, time.Time{}, time.Time{}, 0, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// количество столов
	if scanner.Scan() {
		numTables, err = strconv.Atoi(scanner.Text())
		if err != nil {
			return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
		}
	}
	if numTables <= 0 {
		return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
	}

	// время начала и окончания работы
	if scanner.Scan() {
		times := strings.Split(scanner.Text(), " ")
		if len(times) != 2 {
			return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
		}

		startTime, err = parseTime(times[0])
		if err != nil {
			return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
		}

		endTime, err = parseTime(times[1])
		if err != nil {
			return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
		}

		if startTime.After(endTime) {
			return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
		}
	}

	// Стоимость часа
	if scanner.Scan() {
		hourCost, err = strconv.Atoi(scanner.Text())
		if err != nil {
			return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
		}
	}
	if hourCost <= 0 {
		return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
	}

	// События
	prevEventTime, _ := time.Parse("15:04", "00:00")
	for scanner.Scan() {
		eventItems := strings.Split(scanner.Text(), " ")
		if len(eventItems) < 3 {
			return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
		}

		//время события
		eventTime, err := parseTime(eventItems[0])
		if err != nil {
			return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
		}
		if eventTime.Before(prevEventTime) {
			return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
		}
		prevEventTime = eventTime

		// идентификатор события
		eventID, err := strconv.Atoi(eventItems[1])
		if err != nil || eventID < 1 || eventID > 4 {
			return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
		}

		// имя клиента
		clientName := eventItems[2]
		err = checkName(clientName)
		if err != nil || eventID < 1 || eventID > 4 {
			return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
		}

		// номер стола
		tableNum := 0
		switch eventID {
		case 2:
			if len(eventItems) != 4 {
				return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
			}
			tableNum, err = strconv.Atoi(eventItems[3])
			if err != nil || tableNum < 1 || tableNum > numTables {
				return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
			}
		default:
			if len(eventItems) != 3 {
				return 0, time.Time{}, time.Time{}, 0, nil, errors.New(scanner.Text())
			}
		}
		events = append(events, Event{Time: eventTime, Id: eventID, ClientName: clientName, TableNum: tableNum})
	}

	if err := scanner.Err(); err != nil {
		return 0, time.Time{}, time.Time{}, 0, nil, err
	}

	return numTables, startTime, endTime, hourCost, events, nil
}
