package ingest

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"io"
	"strings"
	"time"
)

type Record struct {
	TurbineID   string    `json:"turbine_id"`
	At          time.Time `json:"at"`
	PowerKW     float64   `json:"power_kw"`
	WindMS      float64   `json:"wind_ms"`
	Temperature float64   `json:"temperature"`
}
type Decoder struct {
	MaxRecordBytes int
	MaxRecords     int
}

func (d Decoder) Decode(ctx context.Context, reader io.Reader) ([]domain.Sample, error) {
	if d.MaxRecordBytes <= 0 {
		d.MaxRecordBytes = 64 * 1024
	}
	if d.MaxRecords <= 0 {
		d.MaxRecords = 10000
	}
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 1024), d.MaxRecordBytes)
	samples := []domain.Sample{}
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var record Record
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return nil, fmt.Errorf("decode telemetry line: %w", err)
		}
		sample := domain.Sample{TurbineID: record.TurbineID, At: record.At, PowerKW: record.PowerKW, WindMS: record.WindMS, Temperature: record.Temperature}
		if err := domain.ValidateSample(sample); err != nil {
			return nil, err
		}
		samples = append(samples, sample)
		if len(samples) > d.MaxRecords {
			return nil, fmt.Errorf("%w: too many records", domain.ErrValidation)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan telemetry: %w", err)
	}
	return samples, nil
}
func Encode(samples []domain.Sample) []byte {
	var builder strings.Builder
	for _, sample := range samples {
		raw, _ := json.Marshal(Record{TurbineID: sample.TurbineID, At: sample.At, PowerKW: sample.PowerKW, WindMS: sample.WindMS, Temperature: sample.Temperature})
		builder.Write(raw)
		builder.WriteByte('\n')
	}
	return []byte(builder.String())
}
