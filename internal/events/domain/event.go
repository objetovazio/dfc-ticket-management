package domain

import (
	"errors"
	"time"
)

var (
	ErrEventNameRequired = errors.New("Event name is required")
	ErrEventDateRequired = errors.New("Event date is required")
	ErrEventCapacityZero = errors.New("Event capacity must be greater than zero")
	ErrEventPriceZero    = errors.New("Event price must be greater than zero")
)

type Rating string

const (
	RatingLivre Rating = "L"
	Rating10    Rating = "L10"
	Rating12    Rating = "L12"
	Rating14    Rating = "L14"
	Rating16    Rating = "L16"
	Rating18    Rating = "L18"
)

type Event struct {
	ID           string
	Name         string
	Location     string
	Organization string
	Rating       Rating
	Date         time.Time
	ImageURL     string
	Capacity     int
	Price        float64
	PartnerID    int
	Spots        []Spot
	Tickets      []Ticket
}

func (e Event) Validate() error {
	if e.Name == "" {
		return ErrEventNameRequired
	}

	if e.Date.Before(time.Now()) {
		return ErrEventDateRequired
	}

	if e.Capacity <= 0 {
		return ErrEventCapacityZero
	}

	if e.Price <= 0 {
		return ErrEventPriceZero
	}

	return nil
}

func (e *Event) AddSpot(name string) (*Spot, error) {
	spot, err := NewSpot(e, name)

	if err != nil {
		return nil, err
	}

	e.Spots = append(e.Spots, *spot)

	return spot, nil
}
