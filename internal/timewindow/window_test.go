package timewindow_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/timewindow"
	"testing"
	"time"
)

func TestOverlapAndCheck(t *testing.T) {
	now := time.Now()
	a := timewindow.Window{ID: "a", TurbineID: "t", Start: now, End: now.Add(time.Hour)}
	b := timewindow.Window{ID: "b", TurbineID: "t", Start: now.Add(30 * time.Minute), End: now.Add(2 * time.Hour)}
	if !timewindow.Overlaps(a, b) {
		t.Fatal("overlap missed")
	}
	if err := timewindow.Check([]timewindow.Window{a, b}); err == nil {
		t.Fatal("overlap accepted")
	}
	if timewindow.Remaining(a, now.Add(2*time.Hour)) != 0 {
		t.Fatal("remaining")
	}
}
