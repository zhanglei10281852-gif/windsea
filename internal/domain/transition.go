package domain

import (
	"fmt"
	"strings"
)

type Transition struct{ Entity, ID, From, To, Actor, Reason string }

func ValidateTransition(entity, from, to string) error {
	if strings.TrimSpace(entity) == "" || strings.TrimSpace(from) == "" || strings.TrimSpace(to) == "" {
		return fmt.Errorf("%w: transition fields", ErrValidation)
	}
	switch entity {
	case "campaign":
		if !CampaignState(from).CanMove(CampaignState(to)) {
			return ErrInvalidState
		}
	case "work_order":
		if !WorkOrderState(from).CanMove(WorkOrderState(to)) {
			return ErrInvalidState
		}
	case "reservation":
		if !ReservationState(from).CanMove(ReservationState(to)) {
			return ErrInvalidState
		}
	case "alert":
		if !AlertState(from).CanMove(AlertState(to)) {
			return ErrInvalidState
		}
	default:
		return fmt.Errorf("%w: unknown entity", ErrValidation)
	}
	return nil
}
func ApplyTransitions(current string, transitions []Transition) (string, error) {
	state := current
	for _, transition := range transitions {
		if transition.From != state {
			return state, fmt.Errorf("%w: expected %s", ErrConflict, state)
		}
		if err := ValidateTransition(transition.Entity, state, transition.To); err != nil {
			return state, err
		}
		state = transition.To
	}
	return state, nil
}
