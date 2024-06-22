package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrSpotNameRequired    = errors.New("Spot name is required")
	ErrInvalidSpotNumber   = errors.New("Spot name must be at least 2 characters long")
	ErrSpotNotFound        = errors.New("Spot not found")
	ErrSpotAlreadyReserved = errors.New("Spot already reserved")
	ErrSpotNameStartLetter = errors.New("Spot name must start with a letter")
	ErrSpotEndNumber       = errors.New("Spot name must end with a number")
)

type SpotStatus string

const (
	SpotStatusAvailable SpotStatus = "available"
	SpotStatusSold      SpotStatus = "sold"
)

type Spot struct {
	ID       string
	EventID  string
	Name     string
	Status   SpotStatus
	TicketID string
}

func (s *Spot) Validate() error {
	if s.Name == "" {
		return ErrSpotNameRequired
	}

	if len(s.Name) < 2 {
		return ErrInvalidSpotNumber
	}

	if s.Name[0] < 'A' || s.Name[0] > 'Z' {
		return ErrSpotNameStartLetter
	}

	if s.Name[1] < '0' || s.Name[1] > '9' {
		return ErrSpotEndNumber
	}

	return nil
}

func NewSpot(event *Event, name string) (*Spot, error) {
	spot := &Spot{
		ID:      uuid.New().String(),
		EventID: event.ID,
		Name:    name,
		Status:  SpotStatusAvailable,
	}

	if err := spot.Validate(); err != nil {
		return nil, err
	}

	return spot, nil
}

func (s *Spot) Reserve(ticketID string) error {
	if s.Status == SpotStatusSold {
		return ErrSpotAlreadyReserved
	}

	s.Status = SpotStatusSold
	s.TicketID = ticketID

	return nil
}
