package usecase

import (
	"github.com/devfullcycle/imersao18/golang/internal/events/domain"
	"github.com/devfullcycle/imersao18/golang/internal/events/domain/infra/service"
)

// Ticket struct
// ID         string
type BuyTicketInputDTO struct {
	EventID    string   `json:"event_id"`
	Spots      []string `json:"spot"`
	TicketType string   `json:"ticket_type"`
	CardHash   string   `json:"card_hash"`
	Email      string   `json:"email"`
}

type BuyTicketsOutputDTO struct {
	Tickets []TicketDTO `json:"tickets"`
}

type BuyTicketsUseCase struct {
	repo           domain.EventRepository
	partnerFactory service.PartnerFactory
}

func NewBuyTicketsUseCase(repo domain.EventRepository, partnerFactory service.PartnerFactory) *BuyTicketsUseCase {
	return &BuyTicketsUseCase{
		repo:           repo,
		partnerFactory: partnerFactory,
	}
}

func (uc *BuyTicketsUseCase) Execute(dto BuyTicketInputDTO) (*BuyTicketsOutputDTO, error) {
	event, err := uc.repo.FindEventById(dto.EventID)
	if err != nil {
		return nil, err
	}

	request := &service.ReservationRequest{
		EventID:    dto.EventID,
		Spots:      dto.Spots,
		TicketType: dto.TicketType,
		CardHash:   dto.CardHash,
		Email:      dto.Email,
	}

	partnerService, err := uc.partnerFactory.CreatePartner(event.PartnerID)
	if err != nil {
		return nil, err
	}

	reservationResponse, err := partnerService.MakeReservation(request)
	if err != nil {
		return nil, err
	}

	tickets := make([]domain.Ticket, len(reservationResponse))
	for i, reservation := range reservationResponse {
		spot, err := uc.repo.FindSpotByName(event.ID, reservation.Spot)
		if err != nil {
			return nil, err
		}

		ticket, err := domain.NewTicket(event, spot, domain.TicketType(reservation.TicketType))
		if err != nil {
			return nil, err
		}

		err = uc.repo.CreateTicket(ticket)
		if err != nil {
			return nil, err
		}

		spot.Reserve(ticket.ID)
		err = uc.repo.ReserveSpot(spot.ID, ticket.ID)
		if err != nil {
			return nil, err
		}

		tickets[i] = *ticket
	}

	ticketsDTOs := make([]TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketsDTOs[i] = TicketDTO{
			ID:         ticket.ID,
			SpotID:     ticket.Spot.ID,
			TicketType: string(ticket.TicketType),
			Price:      ticket.Price,
		}
	}

	return &BuyTicketsOutputDTO{
		Tickets: ticketsDTOs,
	}, nil
}
