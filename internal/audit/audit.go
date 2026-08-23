package audit

import (
	"encoding/json"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"time"
)

type Recorder interface{ Append(domain.AuditEvent) error }
type EventBuilder struct{ Now func() time.Time }

func (b EventBuilder) Build(id, farm, actor, object, objectID, action, result, request string, metadata map[string]any) (domain.AuditEvent, error) {
	if id == "" || farm == "" || object == "" || objectID == "" || action == "" {
		return domain.AuditEvent{}, domain.ErrValidation
	}
	raw, err := json.Marshal(metadata)
	if err != nil {
		return domain.AuditEvent{}, fmt.Errorf("metadata: %w", err)
	}
	return domain.AuditEvent{ID: id, FarmID: farm, ActorID: actor, ObjectType: object, ObjectID: objectID, Action: action, Result: result, RequestID: request, Metadata: string(raw), CreatedAt: b.Now()}, nil
}
func DecodeMetadata(event domain.AuditEvent) (map[string]any, error) {
	value := map[string]any{}
	if event.Metadata == "" {
		return value, nil
	}
	if err := json.Unmarshal([]byte(event.Metadata), &value); err != nil {
		return nil, fmt.Errorf("decode metadata: %w", err)
	}
	return value, nil
}
