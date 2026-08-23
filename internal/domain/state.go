package domain

type CampaignState string

const (
	CampaignDraft     CampaignState = "draft"
	CampaignReview    CampaignState = "review"
	CampaignPublished CampaignState = "published"
	CampaignClosed    CampaignState = "closed"
)

func (s CampaignState) CanMove(next CampaignState) bool {
	switch s {
	case CampaignDraft:
		return next == CampaignReview
	case CampaignReview:
		return next == CampaignPublished || next == CampaignDraft
	case CampaignPublished:
		return next == CampaignClosed
	default:
		return false
	}
}

type WorkOrderState string

const (
	WorkQueued     WorkOrderState = "queued"
	WorkAssigned   WorkOrderState = "assigned"
	WorkInProgress WorkOrderState = "in_progress"
	WorkBlocked    WorkOrderState = "blocked"
	WorkCompleted  WorkOrderState = "completed"
	WorkCancelled  WorkOrderState = "cancelled"
)

func (s WorkOrderState) CanMove(next WorkOrderState) bool {
	switch s {
	case WorkQueued:
		return next == WorkAssigned || next == WorkCancelled
	case WorkAssigned:
		return next == WorkInProgress || next == WorkBlocked || next == WorkCancelled
	case WorkInProgress:
		return next == WorkBlocked || next == WorkCompleted || next == WorkCancelled
	case WorkBlocked:
		return next == WorkInProgress || next == WorkCancelled
	default:
		return false
	}
}

type ReservationState string

const (
	ReservationHeld     ReservationState = "held"
	ReservationConsumed ReservationState = "consumed"
	ReservationReleased ReservationState = "released"
)

func (s ReservationState) CanMove(next ReservationState) bool {
	return (s == ReservationHeld && (next == ReservationConsumed || next == ReservationReleased))
}

type AlertState string

const (
	AlertOpen         AlertState = "open"
	AlertAcknowledged AlertState = "acknowledged"
	AlertResolved     AlertState = "resolved"
)

func (s AlertState) CanMove(next AlertState) bool {
	return (s == AlertOpen && next == AlertAcknowledged) || (s == AlertAcknowledged && next == AlertResolved)
}
