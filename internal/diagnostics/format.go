package diagnostics

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

type TimelineEvent struct {
	At       time.Time `json:"at"`
	Kind     string    `json:"kind"`
	Message  string    `json:"message"`
	ObjectID string    `json:"object_id"`
}

func NormalizeTimeline(events []TimelineEvent) []TimelineEvent {
	copy := append([]TimelineEvent(nil), events...)
	sort.SliceStable(copy, func(i, j int) bool {
		if copy[i].At.Equal(copy[j].At) {
			return copy[i].Kind < copy[j].Kind
		}
		return copy[i].At.Before(copy[j].At)
	})
	return copy
}
func FilterTimeline(events []TimelineEvent, objectID string) []TimelineEvent {
	result := []TimelineEvent{}
	for _, event := range events {
		if objectID == "" || event.ObjectID == objectID {
			result = append(result, event)
		}
	}
	return result
}
func WindowTimeline(events []TimelineEvent, start, end time.Time) []TimelineEvent {
	result := []TimelineEvent{}
	for _, event := range events {
		if !event.At.Before(start) && event.At.Before(end) {
			result = append(result, event)
		}
	}
	return NormalizeTimeline(result)
}
func EncodeTimeline(events []TimelineEvent) ([]byte, error) {
	if events == nil {
		events = []TimelineEvent{}
	}
	return json.Marshal(NormalizeTimeline(events))
}
func DecodeTimeline(raw []byte) ([]TimelineEvent, error) {
	var events []TimelineEvent
	if err := json.Unmarshal(raw, &events); err != nil {
		return nil, fmt.Errorf("decode timeline: %w", err)
	}
	return NormalizeTimeline(events), nil
}
func DescribeTimeline(events []TimelineEvent) string {
	events = NormalizeTimeline(events)
	parts := make([]string, 0, len(events))
	for _, event := range events {
		message := strings.TrimSpace(event.Message)
		if message == "" {
			message = "(no message)"
		}
		parts = append(parts, fmt.Sprintf("%s %s: %s", event.At.UTC().Format(time.RFC3339), event.Kind, message))
	}
	return strings.Join(parts, "\n")
}
func GroupByKind(events []TimelineEvent) map[string][]TimelineEvent {
	result := map[string][]TimelineEvent{}
	for _, event := range NormalizeTimeline(events) {
		result[event.Kind] = append(result[event.Kind], event)
	}
	return result
}
func LatestByObject(events []TimelineEvent) map[string]TimelineEvent {
	result := map[string]TimelineEvent{}
	for _, event := range NormalizeTimeline(events) {
		current, ok := result[event.ObjectID]
		if !ok || event.At.After(current.At) {
			result[event.ObjectID] = event
		}
	}
	return result
}
func HasFailure(events []TimelineEvent) bool {
	for _, event := range events {
		if strings.Contains(strings.ToLower(event.Kind), "fail") || strings.Contains(strings.ToLower(event.Message), "error") {
			return true
		}
	}
	return false
}
func Count(events []TimelineEvent, kind string) int {
	count := 0
	for _, event := range events {
		if event.Kind == kind {
			count++
		}
	}
	return count
}
func Merge(groups ...[]TimelineEvent) []TimelineEvent {
	size := 0
	for _, group := range groups {
		size += len(group)
	}
	merged := make([]TimelineEvent, 0, size)
	for _, group := range groups {
		merged = append(merged, group...)
	}
	return NormalizeTimeline(merged)
}
func Trim(events []TimelineEvent, limit int) []TimelineEvent {
	if limit < 0 {
		limit = 0
	}
	events = NormalizeTimeline(events)
	if len(events) <= limit {
		return events
	}
	return append([]TimelineEvent(nil), events[len(events)-limit:]...)
}
func ValidateTimeline(events []TimelineEvent) error {
	for _, event := range events {
		if event.At.IsZero() || strings.TrimSpace(event.Kind) == "" || strings.TrimSpace(event.ObjectID) == "" {
			return fmt.Errorf("invalid timeline event")
		}
	}
	return nil
}
