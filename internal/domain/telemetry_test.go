package domain_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"testing"
	"time"
)

func TestTelemetryValidation(t *testing.T) {
	now := time.Now()
	valid := domain.Sample{TurbineID: "t1", At: now, PowerKW: 100, WindMS: 8, Temperature: 12}
	if err := domain.ValidateSample(valid); err != nil {
		t.Fatal(err)
	}
	for _, sample := range []domain.Sample{{}, {TurbineID: "t1", At: now, PowerKW: -1}, {TurbineID: "t1", At: now, Temperature: 100}} {
		if err := domain.ValidateSample(sample); err == nil {
			t.Errorf("expected invalid sample %#v", sample)
		}
	}
}
func TestTelemetrySummary(t *testing.T) {
	now := time.Now()
	summary := domain.Summarize([]domain.Sample{{TurbineID: "t1", At: now, PowerKW: 100, WindMS: 5, Temperature: 10}, {TurbineID: "t1", At: now.Add(time.Minute), PowerKW: 200, WindMS: 7, Temperature: 11}, {TurbineID: "t1", At: now, PowerKW: -1}})
	if summary.Accepted != 2 || summary.Rejected != 1 || summary.PeakPower != 200 {
		t.Fatalf("summary=%+v", summary)
	}
	if summary.MeanPower != 150 || summary.MeanWind != 6 {
		t.Fatalf("means=%+v", summary)
	}
}
