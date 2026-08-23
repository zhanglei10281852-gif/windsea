package domain

import "time"

type User struct {
	ID, FarmID, Email, Name, Role string
	Active                        bool
	CreatedAt                     time.Time
}
type Session struct {
	ID, UserID, TokenHash string
	ExpiresAt, RevokedAt  *time.Time
}
type Farm struct {
	ID, Name, Timezone string
	TurbineCount       int
	CreatedAt          time.Time
}
type Turbine struct {
	ID, FarmID, Name, Model string
	RatedKW                 int
	Status                  string
	Version                 int
	CreatedAt               time.Time
}
type Campaign struct {
	ID, FarmID, Name, State string
	StartAt, EndAt          time.Time
	Version                 int
	CreatedAt, UpdatedAt    time.Time
}
type Inspection struct {
	ID, CampaignID, TurbineID, InspectorID, Status, Notes string
	CompletedAt                                           *time.Time
	CreatedAt                                             time.Time
}
type WorkOrder struct {
	ID, FarmID, CampaignID, TurbineID, AssigneeID, Title, State string
	Priority, Version                                           int
	DueAt                                                       *time.Time
	CreatedAt, UpdatedAt                                        time.Time
}
type Part struct {
	ID, FarmID, SKU, Description string
	OnHand, Reserved             int
	Version                      int
	CreatedAt                    time.Time
}
type Reservation struct {
	ID, PartID, WorkOrderID, RequestedBy, State string
	Quantity, Version                           int
	CreatedAt, UpdatedAt                        time.Time
}
type Alert struct {
	ID, FarmID, TurbineID, Code, Severity, State, Message, DedupKey string
	OccurredAt, AcknowledgedAt, ResolvedAt                          *time.Time
	CreatedAt                                                       time.Time
}
type TelemetryBatch struct {
	ID, FarmID, Source, State   string
	Samples, Accepted, Rejected int
	ReceivedAt, CompletedAt     *time.Time
}
type Handoff struct {
	ID, WorkOrderID, ContractorID, State, Notes string
	AcceptedAt, CompletedAt                     *time.Time
	CreatedAt                                   time.Time
}
type MaintenanceWindow struct {
	ID, FarmID, TurbineID, Name string
	StartsAt, EndsAt            time.Time
	State                       string
}
type AuditEvent struct {
	ID, FarmID, ActorID, ObjectType, ObjectID, Action, Result, RequestID string
	Metadata                                                             string
	CreatedAt                                                            time.Time
}
type Page[T any] struct {
	Items         []T
	Total         int
	Limit, Offset int
}
