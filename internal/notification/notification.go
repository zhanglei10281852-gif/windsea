package notification

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"sync"
	"time"
)

type Message struct {
	ID, FarmID, Recipient, Subject, Body string
	CreatedAt                            time.Time
}
type Sender interface {
	Send(context.Context, Message) error
}
type MemorySender struct {
	mu       sync.Mutex
	Messages []Message
	Failures int
}

func (s *MemorySender) Send(ctx context.Context, message Message) error {
	if err := ctx.Err(); err != nil {
		return errors.Join(domain.ErrCancelled, err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Failures > 0 {
		s.Failures--
		return errors.New("temporary notification failure")
	}
	s.Messages = append(s.Messages, message)
	return nil
}
func BuildAlertMessage(alert domain.Alert, recipient string, at time.Time) (Message, error) {
	if alert.ID == "" || recipient == "" {
		return Message{}, domain.ErrValidation
	}
	return Message{ID: "alert-" + alert.ID, FarmID: alert.FarmID, Recipient: recipient, Subject: fmt.Sprintf("Wind turbine alert: %s", alert.Code), Body: alert.Message, CreatedAt: at}, nil
}
func Deliver(ctx context.Context, sender Sender, message Message, attempts int) error {
	if attempts < 1 {
		attempts = 1
	}
	var last error
	for i := 0; i < attempts; i++ {
		if err := sender.Send(ctx, message); err == nil {
			return nil
		} else {
			last = err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(i+1) * time.Millisecond):
		}
	}
	return last
}
