package forecast_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/forecast"
	"testing"
	"time"
)

func TestForecastFunctions(t *testing.T) {
	now := time.Now()
	points := []forecast.Point{{At: now, Value: 1}, {At: now.Add(time.Minute), Value: 3}, {At: now.Add(2 * time.Minute), Value: 5}}
	average := forecast.MovingAverage(points, 2)
	if average[1].Value != 2 || average[2].Value != 4 {
		t.Fatalf("average=%+v", average)
	}
	if forecast.Trend(points) != 2 {
		t.Fatal("trend")
	}
	clamped := forecast.Clamp(points, 2, 4)
	if clamped[0].Value != 2 || clamped[2].Value != 4 {
		t.Fatal("clamp")
	}
}
