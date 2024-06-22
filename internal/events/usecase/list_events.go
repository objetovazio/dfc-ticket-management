package usecase

import "github.com/devfullcycle/imersao18/golang/internal/events/domain"

type ListEventsOutputDTO struct {
	Events []EventDTO `json:"events"`
}

type ListEventsUseCase struct {
	repo domain.EventRepository
}

func NewListEventsUseCase(repo domain.EventRepository) *ListEventsUseCase {
	return &ListEventsUseCase{repo: repo}
}

func (us *ListEventsUseCase) Execute() (*ListEventsOutputDTO, error) {
	events, err := us.repo.ListEvents()
	if err != nil {
		return nil, err
	}

	eventsDTOs := make([]EventDTO, len(events))
	for i, event := range events {
		eventsDTOs[i] = EventDTO{
			ID:           event.ID,
			Name:         event.Name,
			Location:     event.Location,
			Organization: event.Organization,
			Rating:       string(event.Rating),
			Date:         event.Date.Format("2006-01-02 15:04:05"),
			ImageURL:     event.ImageURL,
			Capacity:     event.Capacity,
			Price:        event.Price,
			PartnerID:    event.PartnerID,
		}
	}

	return &ListEventsOutputDTO{Events: eventsDTOs}, nil
}
