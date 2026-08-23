package repository

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
)

func (r *Repository) CreatePart(ctx context.Context, part domain.Part) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO parts(id,farm_id,sku,description,on_hand,reserved,version,created_at) VALUES(?,?,?,?,?,?,?,?)", part.ID, part.FarmID, part.SKU, part.Description, part.OnHand, part.Reserved, part.Version, text(part.CreatedAt))
	if err != nil {
		return fmt.Errorf("create part: %w", err)
	}
	return nil
}
func (r *Repository) GetPart(ctx context.Context, id string) (domain.Part, error) {
	var p domain.Part
	var created string
	err := r.DB.QueryRowContext(ctx, "SELECT id,farm_id,sku,description,on_hand,reserved,version,created_at FROM parts WHERE id=?", id).Scan(&p.ID, &p.FarmID, &p.SKU, &p.Description, &p.OnHand, &p.Reserved, &p.Version, &created)
	if err != nil {
		if isNoRows(err) {
			return p, domain.ErrNotFound
		}
		return p, err
	}
	p.CreatedAt, _ = parse(created)
	return p, nil
}
func (r *Repository) HoldReservation(ctx context.Context, res domain.Reservation) error {
	return r.WithReservationTx(ctx, res, domain.ReservationHeld)
}
func (r *Repository) WithReservationTx(ctx context.Context, res domain.Reservation, state domain.ReservationState) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin reservation: %w", err)
	}
	defer tx.Rollback()
	var onHand, reserved, version int
	if err := tx.QueryRowContext(ctx, "SELECT on_hand,reserved,version FROM parts WHERE id=?", res.PartID).Scan(&onHand, &reserved, &version); err != nil {
		return fmt.Errorf("load part: %w", err)
	}
	available := onHand - reserved
	if available < 0 {
		return domain.ErrCapacity
	}
	if res.Quantity <= 0 {
		return domain.ErrValidation
	}
	if _, err := tx.ExecContext(ctx, "UPDATE parts SET reserved=reserved+?,version=version+1 WHERE id=? AND version=?", res.Quantity, res.PartID, version); err != nil {
		return fmt.Errorf("reserve part: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO reservations(id,part_id,work_order_id,requested_by,state,quantity,version,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)", res.ID, res.PartID, res.WorkOrderID, res.RequestedBy, string(state), res.Quantity, res.Version, text(res.CreatedAt), text(res.UpdatedAt)); err != nil {
		return fmt.Errorf("create reservation: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit reservation: %w", err)
	}
	return nil
}
func (r *Repository) MoveReservation(ctx context.Context, id, from, to string, at string) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var partID string
	var quantity int
	var version int
	if err := tx.QueryRowContext(ctx, "SELECT part_id,quantity,version FROM reservations WHERE id=? AND state=?", id, from).Scan(&partID, &quantity, &version); err != nil {
		if isNoRows(err) {
			return domain.ErrConflict
		}
		return err
	}
	if _, err := tx.ExecContext(ctx, "UPDATE reservations SET state=?,version=version+1,updated_at=? WHERE id=? AND state=? AND version=?", to, at, id, from, version); err != nil {
		return err
	}
	if to == string(domain.ReservationConsumed) || to == string(domain.ReservationReleased) {
		if _, err := tx.ExecContext(ctx, "UPDATE parts SET reserved=reserved-?,version=version+1 WHERE id=? AND reserved>=?", quantity, partID, quantity); err != nil {
			return err
		}
	}
	return tx.Commit()
}
