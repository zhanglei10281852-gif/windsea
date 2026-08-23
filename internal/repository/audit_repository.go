package repository

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
)

func (r *Repository) AppendAudit(ctx context.Context, event domain.AuditEvent) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO audit_events(id,farm_id,actor_id,object_type,object_id,action,result,request_id,metadata,created_at) VALUES(?,?,?,?,?,?,?,?,?,?)", event.ID, event.FarmID, event.ActorID, event.ObjectType, event.ObjectID, event.Action, event.Result, event.RequestID, event.Metadata, text(event.CreatedAt))
	if err != nil {
		return fmt.Errorf("append audit: %w", err)
	}
	return nil
}

func (r *Repository) AuditFailureMessage(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("audit persistence failed after business commit: %v", err)
}
func (r *Repository) ListAudit(ctx context.Context, farmID, objectID string) ([]domain.AuditEvent, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id,farm_id,COALESCE(actor_id,''),object_type,object_id,action,result,request_id,metadata,created_at FROM audit_events WHERE farm_id=? AND object_id=? ORDER BY created_at", farmID, objectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.AuditEvent{}
	for rows.Next() {
		var e domain.AuditEvent
		var created string
		if err := rows.Scan(&e.ID, &e.FarmID, &e.ActorID, &e.ObjectType, &e.ObjectID, &e.Action, &e.Result, &e.RequestID, &e.Metadata, &created); err != nil {
			return nil, err
		}
		e.CreatedAt, _ = parse(created)
		items = append(items, e)
	}
	return items, rows.Err()
}
