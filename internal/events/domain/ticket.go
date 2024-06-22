package domain

import (
	"errors"

	"github.com/google/uuid"
)

type TicketType string

var (
	ErrTicketPriceZero = errors.New("ticket price cannot be zero")
)

const (
	TicketTypeHalf TicketType = "half"
	TicketTypeFull TicketType = "full"
)

type Ticket struct {
	ID         string
	EventID    string
	Spot       *Spot
	TicketType TicketType
	Price      float64
}

func IsValidTicketType(tycketType TicketType) bool {
	return tycketType == TicketTypeHalf || tycketType == TicketTypeFull
}

func (t *Ticket) CalculatePrice() {
	if t.TicketType == TicketTypeHalf {
		t.Price /= 2
	}
}

func (t *Ticket) Validate() error {
	// Validate Price
	if t.Price < 0 {
		return ErrTicketPriceZero
	}

	return nil
}

func NewTicket(event *Event, spot *Spot, ticketType TicketType) (*Ticket, error) {
	if !IsValidTicketType(ticketType) {
		return nil, errors.New("invalid ticket type")
	}

	ticket := &Ticket{
		ID:         uuid.New().String(),
		EventID:    event.ID,
		Spot:       spot,
		TicketType: ticketType,
		Price:      event.Price,
	}

	ticket.CalculatePrice()
	if err := ticket.Validate(); err != nil {
		return nil, err
	}

	return ticket, nil
}
