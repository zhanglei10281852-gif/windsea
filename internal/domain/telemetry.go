package domain

import (
	"fmt"
	"math"
	"time"
)

type Sample struct {
	TurbineID                    string
	At                           time.Time
	PowerKW, WindMS, Temperature float64
}
type TelemetrySummary struct {
	Samples, Accepted, Rejected    int
	MeanPower, PeakPower, MeanWind float64
	WindowStart, WindowEnd         time.Time
}

func ValidateSample(sample Sample) error {
	if sample.TurbineID == "" || sample.At.IsZero() {
		return ErrValidation
	}
	if sample.PowerKW < 0 || sample.WindMS < 0 || sample.Temperature < -80 || sample.Temperature > 90 {
		return fmt.Errorf("%w: sample range", ErrValidation)
	}
	return nil
}
func Summarize(samples []Sample) TelemetrySummary {
	summary := TelemetrySummary{Samples: len(samples)}
	for _, sample := range samples {
		if ValidateSample(sample) != nil {
			summary.Rejected++
			continue
		}
		summary.Accepted++
		summary.MeanPower += sample.PowerKW
		summary.MeanWind += sample.WindMS
		if sample.PowerKW > summary.PeakPower {
			summary.PeakPower = sample.PowerKW
		}
		if summary.WindowStart.IsZero() || sample.At.Before(summary.WindowStart) {
			summary.WindowStart = sample.At
		}
		if sample.At.After(summary.WindowEnd) {
			summary.WindowEnd = sample.At
		}
	}
	if summary.Accepted > 0 {
		summary.MeanPower /= float64(summary.Accepted)
		summary.MeanWind /= float64(summary.Accepted)
	}
	if math.IsNaN(summary.MeanPower) {
		summary.MeanPower = 0
	}
	return summary
}
