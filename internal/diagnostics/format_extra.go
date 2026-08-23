package diagnostics

import (
	"fmt"
	"strings"
	"time"
)

func Earliest(events []TimelineEvent) (TimelineEvent, bool) {
	if len(events) == 0 {
		return TimelineEvent{}, false
	}
	earliest := events[0]
	for _, event := range events[1:] {
		if event.At.Before(earliest.At) {
			earliest = event
		}
	}
	return earliest, true
}
func Latest(events []TimelineEvent) (TimelineEvent, bool) {
	if len(events) == 0 {
		return TimelineEvent{}, false
	}
	latest := events[0]
	for _, event := range events[1:] {
		if event.At.After(latest.At) {
			latest = event
		}
	}
	return latest, true
}
func Duration(events []TimelineEvent) time.Duration {
	first, ok := Earliest(events)
	if !ok {
		return 0
	}
	last, ok := Latest(events)
	if !ok {
		return 0
	}
	return last.At.Sub(first.At)
}
func Kinds(events []TimelineEvent) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, event := range events {
		if !seen[event.Kind] {
			seen[event.Kind] = true
			result = append(result, event.Kind)
		}
	}
	return result
}
func Objects(events []TimelineEvent) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, event := range events {
		if !seen[event.ObjectID] {
			seen[event.ObjectID] = true
			result = append(result, event.ObjectID)
		}
	}
	return result
}
func Contains(events []TimelineEvent, needle string) bool {
	needle = strings.ToLower(needle)
	for _, event := range events {
		if strings.Contains(strings.ToLower(event.Message), needle) {
			return true
		}
	}
	return false
}
func RequireKind(events []TimelineEvent, kind string) error {
	for _, event := range events {
		if event.Kind == kind {
			return nil
		}
	}
	return fmt.Errorf("missing event kind %s", kind)
}
func CountObject(events []TimelineEvent, objectID string) int {
	count := 0
	for _, event := range events {
		if event.ObjectID == objectID {
			count++
		}
	}
	return count
}
func Copy(events []TimelineEvent) []TimelineEvent { return append([]TimelineEvent(nil), events...) }
func Before(events []TimelineEvent, at time.Time) []TimelineEvent {
	result := []TimelineEvent{}
	for _, event := range events {
		if event.At.Before(at) {
			result = append(result, event)
		}
	}
	return result
}
func After(events []TimelineEvent, at time.Time) []TimelineEvent {
	result := []TimelineEvent{}
	for _, event := range events {
		if event.At.After(at) {
			result = append(result, event)
		}
	}
	return result
}
func KindsString(events []TimelineEvent) string  { return strings.Join(Kinds(events), ",") }
func ObjectString(events []TimelineEvent) string { return strings.Join(Objects(events), ",") }
func Empty(events []TimelineEvent) bool          { return len(events) == 0 }
func LimitKinds(events []TimelineEvent, limit int) []TimelineEvent {
	if limit <= 0 {
		return []TimelineEvent{}
	}
	result := []TimelineEvent{}
	allowed := map[string]bool{}
	for _, kind := range Kinds(events) {
		if len(allowed) >= limit {
			break
		}
		allowed[kind] = true
	}
	for _, event := range events {
		if allowed[event.Kind] {
			result = append(result, event)
		}
	}
	return result
}
func Summaries(events []TimelineEvent) map[string]string {
	result := map[string]string{}
	for _, event := range events {
		if _, ok := result[event.ObjectID]; !ok {
			result[event.ObjectID] = event.Message
		}
	}
	return result
}
