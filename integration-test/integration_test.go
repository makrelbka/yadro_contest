package integration_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"yadro/config"
	"yadro/internal/controller"
	"yadro/internal/usecase/biathlon"
	"yadro/internal/usecase/repository"
)

func TestFromReadMe(t *testing.T) {
	configContent := `{
        "laps": 2,
        "lapLen": 3651,
        "penaltyLen": 50,
        "firingLines": 1,
        "start": "09:30:00",
        "startDelta": "00:00:30"
    }`
	configFile := "test_config.json"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err, "failed to create config file")
	defer os.Remove(configFile)

	eventsContent := `[09:05:59.867] 1 1
[09:15:00.841] 2 1 09:30:00.000
[09:29:45.734] 3 1
[09:30:01.005] 4 1
[09:49:31.659] 5 1 1
[09:49:33.123] 6 1 1
[09:49:34.650] 6 1 2
[09:49:35.937] 6 1 4
[09:49:37.364] 6 1 5
[09:49:38.339] 7 1
[09:49:55.915] 8 1
[09:51:48.391] 9 1
[09:59:03.872] 10 1
[09:59:03.872] 11 1 Lost in the forest`
	eventsFile := "test_events.txt"
	err = os.WriteFile(eventsFile, []byte(eventsContent), 0644)
	require.NoError(t, err, "failed to create events file")
	defer os.Remove(eventsFile)

	cfg, err := config.Load(configFile)
	require.NoError(t, err, "failed to load config")

	repo := inmemory.NewInMemoryRepository(cfg)
	service := biathlon.NewCompetitorService(repo)
	processor := controller.NewEventProcessor(service)

	var logs []string
	err = controller.ProcessEvents(eventsFile, &logs, processor)
	require.NoError(t, err, "failed to process events")

	report := controller.GenerateFinalReport(cfg, repo.Competitors)

	expectedLogs := []string{
		"[09:05:59.867] The competitor(1) registered",
		"[09:15:00.841] The start time for the competitor(1) was set by a draw to 09:30:00.000",
		"[09:29:45.734] The competitor(1) is on the start line",
		"[09:30:01.005] The competitor(1) has started",
		"[09:49:31.659] The competitor(1) is on the firing range(1)",
		"[09:49:33.123] The target(1) has been hit by competitor(1)",
		"[09:49:34.650] The target(2) has been hit by competitor(1)",
		"[09:49:35.937] The target(4) has been hit by competitor(1)",
		"[09:49:37.364] The target(5) has been hit by competitor(1)",
		"[09:49:38.339] The competitor(1) left the firing range",
		"[09:49:55.915] The competitor(1) entered the penalty laps",
		"[09:51:48.391] The competitor(1) left the penalty laps",
		"[09:59:03.872] The competitor(1) ended the main lap",
		"[09:59:03.872] The competitor(1) can`t continue: Lost in the forest",
	}

	expectedReport := []string{
		"[NotFinished] 1 [{00:29:02.867, 2.095}, {,}] {00:01:52.476, 0.445} 4/5",
	}

	require.Equal(t, len(expectedLogs), len(logs), "log count mismatch")
	for i, log := range logs {
		require.Equal(t, expectedLogs[i], log, "unexpected log at index %d", i)
	}

	require.Equal(t, len(expectedReport), len(report), "report line count mismatch")
	for i, line := range report {
		require.Equal(t, expectedReport[i], line, "unexpected report line at index %d", i)
	}
}

func TestTwoCompetitorsSuccessfulFinish(t *testing.T) {
	configContent := `{
        "laps": 2,
        "lapLen": 3651,
        "penaltyLen": 50,
        "firingLines": 1,
        "start": "09:30:00",
        "startDelta": "00:00:30"
    }`
	configFile := "test_config_two_success.json"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err, "failed to create config file")
	defer os.Remove(configFile)

	eventsContent := `[09:05:59.867] 1 1
    [09:06:00.000] 1 2
    [09:15:00.841] 2 1 09:30:00.000
    [09:15:30.841] 2 2 09:30:30.000
    [09:29:45.734] 3 1
    [09:30:15.734] 3 2
    [09:30:01.005] 4 1
    [09:30:31.005] 4 2
    [09:49:31.659] 5 1 1
    [09:49:33.123] 6 1 1
    [09:49:34.650] 6 1 2
    [09:49:35.937] 6 1 4
    [09:49:37.364] 6 1 5
    [09:49:38.339] 7 1
    [09:49:55.915] 8 1
    [09:51:48.391] 9 1
    [09:50:31.659] 5 2 1
    [09:50:33.123] 6 2 1
    [09:50:34.650] 6 2 2
    [09:50:35.937] 6 2 3
    [09:50:37.364] 6 2 4
    [09:50:37.364] 6 2 5
    [09:50:38.339] 7 2
    [09:58:30.000] 10 2
    [09:59:03.872] 10 1
    [10:28:30.000] 5 1 1
    [10:28:33.123] 6 1 1
    [10:28:34.650] 6 1 2
    [10:28:35.937] 6 1 4
    [10:28:37.364] 6 1 5
    [10:28:38.339] 7 1
    [10:28:55.915] 8 1
    [10:32:48.391] 9 1
    [10:29:30.000] 5 2 1
    [10:29:33.123] 6 2 1
    [10:29:34.650] 6 2 2
    [10:29:35.937] 6 2 3
    [10:29:37.364] 6 2 4
    [10:29:37.364] 6 2 5
    [10:29:38.339] 7 2
    [10:58:30.000] 10 2
    [10:59:03.872] 10 1`
	eventsFile := "test_events_two_success.txt"
	err = os.WriteFile(eventsFile, []byte(eventsContent), 0644)
	require.NoError(t, err, "failed to create events file")
	defer os.Remove(eventsFile)

	cfg, err := config.Load(configFile)
	require.NoError(t, err, "failed to load config")

	repo := inmemory.NewInMemoryRepository(cfg)
	service := biathlon.NewCompetitorService(repo)
	processor := controller.NewEventProcessor(service)

	var logs []string
	err = controller.ProcessEvents(eventsFile, &logs, processor)
	require.NoError(t, err, "failed to process events")

	report := controller.GenerateFinalReport(cfg, repo.Competitors)

	expectedLogs := []string{
		"[09:05:59.867] The competitor(1) registered",
		"[09:06:00.000] The competitor(2) registered",
		"[09:15:00.841] The start time for the competitor(1) was set by a draw to 09:30:00.000",
		"[09:15:30.841] The start time for the competitor(2) was set by a draw to 09:30:30.000",
		"[09:29:45.734] The competitor(1) is on the start line",
		"[09:30:15.734] The competitor(2) is on the start line",
		"[09:30:01.005] The competitor(1) has started",
		"[09:30:31.005] The competitor(2) has started",
		"[09:49:31.659] The competitor(1) is on the firing range(1)",
		"[09:49:33.123] The target(1) has been hit by competitor(1)",
		"[09:49:34.650] The target(2) has been hit by competitor(1)",
		"[09:49:35.937] The target(4) has been hit by competitor(1)",
		"[09:49:37.364] The target(5) has been hit by competitor(1)",
		"[09:49:38.339] The competitor(1) left the firing range",
		"[09:49:55.915] The competitor(1) entered the penalty laps",
		"[09:51:48.391] The competitor(1) left the penalty laps",
		"[09:50:31.659] The competitor(2) is on the firing range(1)",
		"[09:50:33.123] The target(1) has been hit by competitor(2)",
		"[09:50:34.650] The target(2) has been hit by competitor(2)",
		"[09:50:35.937] The target(3) has been hit by competitor(2)",
		"[09:50:37.364] The target(4) has been hit by competitor(2)",
		"[09:50:37.364] The target(5) has been hit by competitor(2)",
		"[09:50:38.339] The competitor(2) left the firing range",
		"[09:58:30.000] The competitor(2) ended the main lap",
		"[09:59:03.872] The competitor(1) ended the main lap",
		"[10:28:30.000] The competitor(1) is on the firing range(1)",
		"[10:28:33.123] The target(1) has been hit by competitor(1)",
		"[10:28:34.650] The target(2) has been hit by competitor(1)",
		"[10:28:35.937] The target(4) has been hit by competitor(1)",
		"[10:28:37.364] The target(5) has been hit by competitor(1)",
		"[10:28:38.339] The competitor(1) left the firing range",
		"[10:28:55.915] The competitor(1) entered the penalty laps",
		"[10:32:48.391] The competitor(1) left the penalty laps",
		"[10:29:30.000] The competitor(2) is on the firing range(1)",
		"[10:29:33.123] The target(1) has been hit by competitor(2)",
		"[10:29:34.650] The target(2) has been hit by competitor(2)",
		"[10:29:35.937] The target(3) has been hit by competitor(2)",
		"[10:29:37.364] The target(4) has been hit by competitor(2)",
		"[10:29:37.364] The target(5) has been hit by competitor(2)",
		"[10:29:38.339] The competitor(2) left the firing range",
		"[10:58:30.000] The competitor(2) ended the main lap",
		"[10:59:03.872] The competitor(1) ended the main lap",
	}

	expectedReport := []string{
		"[Finished] 2 [{00:27:58.995, 2.175}, {01:00:00.000, 1.014}] {00:00:00.000, 0.000} 10/10",
		"[Finished] 1 [{00:29:02.867, 2.095}, {01:00:00.000, 1.014}] {00:05:44.952, 0.290} 8/10",
	}

	require.Equal(t, len(expectedLogs), len(logs), "log count mismatch")
	for i, log := range logs {
		require.Equal(t, expectedLogs[i], log, "unexpected log at index %d", i)
	}

	require.Equal(t, len(expectedReport), len(report), "report line count mismatch")
	for i, line := range report {
		require.Equal(t, expectedReport[i], line, "unexpected report line at index %d", i)
	}
}
