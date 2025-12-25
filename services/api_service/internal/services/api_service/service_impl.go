package api_service

import "time"

type Service struct {
	storage   Storage
	eventsPub EventsPublisher
	results   ResultsStore

	ackCoordinator AckCoordinator
	ackNotifier    AckNotifier
}

func New(storage Storage, eventsPub EventsPublisher, results ResultsStore, ackReg AckWaitRegistry, ackWaitTimeout time.Duration) *Service {
	var notifier AckNotifier
	if n, ok := any(ackReg).(AckNotifier); ok {
		notifier = n
	}
	return &Service{
		storage:        storage,
		eventsPub:      eventsPub,
		results:        results,
		ackCoordinator: NewAckCoordinator(ackReg, ackWaitTimeout),
		ackNotifier:    notifier,
	}
}

func (s *Service) Notify(orderID string) {
	if s.ackNotifier == nil {
		return
	}
	s.ackNotifier.Notify(orderID)
}
