package diagnostics

import (
	"fmt"
	"strings"
	"time"
)

func FormatResult(result Result) string {
	state := "ok"
	if !result.OK {
		state = "failed"
	}
	return fmt.Sprintf("%s:%s:%s", result.Name, state, result.Error)
}
func FormatResults(results []Result) string {
	parts := make([]string, 0, len(results))
	for _, result := range results {
		parts = append(parts, FormatResult(result))
	}
	return strings.Join(parts, ";")
}
func ResultAt(result Result, at time.Time) TimelineEvent {
	message := FormatResult(result)
	return TimelineEvent{At: at, Kind: "health", Message: message, ObjectID: result.Name}
}
func ResultsTimeline(results []Result, at time.Time) []TimelineEvent {
	events := make([]TimelineEvent, 0, len(results))
	for _, result := range results {
		events = append(events, ResultAt(result, at))
	}
	return events
}
func Failed(results []Result) []Result {
	failed := []Result{}
	for _, result := range results {
		if !result.OK {
			failed = append(failed, result)
		}
	}
	return failed
}
func Passed(results []Result) []Result {
	passed := []Result{}
	for _, result := range results {
		if result.OK {
			passed = append(passed, result)
		}
	}
	return passed
}
func AverageDuration(results []Result) time.Duration {
	if len(results) == 0 {
		return 0
	}
	total := time.Duration(0)
	for _, result := range results {
		total += result.Duration
	}
	return total / time.Duration(len(results))
}
func Longest(results []Result) Result {
	longest := Result{}
	for _, result := range results {
		if result.Duration > longest.Duration {
			longest = result
		}
	}
	return longest
}
func Shortest(results []Result) Result {
	if len(results) == 0 {
		return Result{}
	}
	shortest := results[0]
	for _, result := range results[1:] {
		if result.Duration < shortest.Duration {
			shortest = result
		}
	}
	return shortest
}
func Names(results []Result) []string {
	names := make([]string, 0, len(results))
	for _, result := range results {
		names = append(names, result.Name)
	}
	return names
}
func AnyFailure(results []Result) bool { return len(Failed(results)) > 0 }
func AllNames(results []Result) string { return strings.Join(Names(results), ",") }
func HasName(results []Result, name string) bool {
	for _, result := range results {
		if result.Name == name {
			return true
		}
	}
	return false
}
func ResultCount(results []Result) int { return len(results) }
