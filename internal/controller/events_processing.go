package controller

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"yadro/internal/entity"
)

func ProcessEvents(filename string, logs *[]string, processor *EventProcessor) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		timeStr := strings.Trim(fields[0], "[]")
		t, err := time.Parse(parseTimeConst, timeStr)
		if err != nil {
			continue
		}

		eventID, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}

		competitorID := fields[2]
		extra := fields[3:]

		processor.Process(entity.Event{
			Time:         t,
			ID:           eventID,
			CompetitorID: competitorID,
			Extra:        extra,
		}, logs)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
