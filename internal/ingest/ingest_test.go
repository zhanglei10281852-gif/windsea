package ingest_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/ingest"
	"strings"
	"testing"
	"time"
)

func TestDecoderRoundTrip(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	raw := ingest.Encode([]domain.Sample{{TurbineID: "t1", At: now, PowerKW: 120, WindMS: 6, Temperature: 10}, {TurbineID: "t2", At: now.Add(time.Minute), PowerKW: 80, WindMS: 5, Temperature: 11}})
	samples, err := (ingest.Decoder{}).Decode(context.Background(), strings.NewReader(string(raw)))
	if err != nil || len(samples) != 2 {
		t.Fatalf("samples=%+v err=%v", samples, err)
	}
}
func TestDecoderRejectsBadRecord(t *testing.T) {
	_, err := (ingest.Decoder{}).Decode(context.Background(), strings.NewReader(`{"turbine_id":"t1","at":"2026-01-01T00:00:00Z","power_kw":-1}`))
	if err == nil {
		t.Fatal("bad record accepted")
	}
}
func TestBatchProcessPartialFailure(t *testing.T) {
	batch := ingest.Batch{ID: "b", FarmID: "f", Source: "scada", Samples: []domain.Sample{{TurbineID: "a", At: time.Now()}, {TurbineID: "b", At: time.Now()}}}
	result := ingest.Process(context.Background(), batch, 2, func(_ context.Context, sample domain.Sample) error {
		if sample.TurbineID == "b" {
			return context.DeadlineExceeded
		}
		return nil
	})
	if result.Accepted != 1 || result.Rejected != 1 {
		t.Fatalf("result=%+v", result)
	}
}
func TestBatchValidation(t *testing.T) {
	if err := ingest.ValidateBatch(ingest.Batch{}); err == nil {
		t.Fatal("empty batch accepted")
	}
}
