package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"time"
)

func PublishCampaignWithAudit(ctx context.Context, db *sql.DB, campaignID string, version int, event domain.AuditEvent, at time.Time) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, "UPDATE campaigns SET state='published',version=version+1,updated_at=? WHERE id=? AND state='review' AND version=?", at.UTC().Format(time.RFC3339Nano), campaignID, version)
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}
	count, _ := result.RowsAffected()
	if count != 1 {
		return domain.ErrConflict
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO audit_events(id,farm_id,actor_id,object_type,object_id,action,result,request_id,metadata,created_at) VALUES(?,?,?,?,?,?,?,?,?,?)", event.ID, event.FarmID, event.ActorID, event.ObjectType, event.ObjectID, event.Action, event.Result, event.RequestID, event.Metadata, at.UTC().Format(time.RFC3339Nano)); err != nil {
		return fmt.Errorf("persist audit: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit campaign: %w", err)
	}
	return nil
}
