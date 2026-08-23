package forecast

import (
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"math"
	"sort"
	"time"
)

type Point struct {
	At    time.Time
	Value float64
}
type Forecast struct {
	TurbineID   string
	Points      []Point
	GeneratedAt time.Time
}

func MovingAverage(points []Point, window int) []Point {
	if window < 1 {
		window = 1
	}
	result := make([]Point, len(points))
	for i := range points {
		start := i - window + 1
		if start < 0 {
			start = 0
		}
		sum := 0.0
		for j := start; j <= i; j++ {
			sum += points[j].Value
		}
		result[i] = Point{At: points[i].At, Value: sum / float64(i-start+1)}
	}
	return result
}
func Clamp(points []Point, min, max float64) []Point {
	result := make([]Point, len(points))
	for i, point := range points {
		value := point.Value
		if value < min {
			value = min
		}
		if value > max {
			value = max
		}
		result[i] = Point{At: point.At, Value: value}
	}
	return result
}
func Trend(points []Point) float64 {
	if len(points) < 2 {
		return 0
	}
	copy := append([]Point(nil), points...)
	sort.Slice(copy, func(i, j int) bool { return copy[i].At.Before(copy[j].At) })
	return (copy[len(copy)-1].Value - copy[0].Value) / float64(len(copy)-1)
}
func ValidPoint(point Point) error {
	if point.At.IsZero() || math.IsNaN(point.Value) || math.IsInf(point.Value, 0) {
		return domain.ErrValidation
	}
	return nil
}
