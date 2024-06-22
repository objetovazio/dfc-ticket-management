package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Partner1 struct {
	BaseURL string
}

type Partner1ReservationRequest struct {
	Spots      []string `json:"spots"`
	TicketKind string   `json:"ticket_kind"`
	Email      string   `json:"email"`
}

type Partner1ReservationResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Spot       string `json:"spot"`
	TicketKind string `json:"ticket_kind"`
	Status     string `json:"status"`
	EventID    string `json:"event_id"`
}

func (p *Partner1) MakeReservation(req *ReservationRequest) ([]ReservationResponse, error) {
	// Instanciate partnerRequest
	partnerRequest := Partner1ReservationRequest{
		Spots:      req.Spots,
		TicketKind: req.TicketType,
		Email:      req.Email,
	}

	// Convert body to json
	body, err := json.Marshal(partnerRequest)
	if err != nil {
		return nil, err
	}

	// Create http call
	url := fmt.Sprintf("%s/events/%s/reserve", p.BaseURL, req.EventID)
	httpRequest, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	// Make call
	client := &http.Client{}
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	// Parse response
	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	// Convert Response
	var partnerResponse []Partner1ReservationResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&partnerResponse); err != nil {
		return nil, err
	}

	// Convert Partner1ReservationResponse to ReservationResponse
	responses := make([]ReservationResponse, len(partnerResponse))
	for i, r := range partnerResponse {
		responses[i] = ReservationResponse{
			ID:     r.ID,
			Spot:   r.Spot,
			Status: r.Status,
		}
	}

	return responses, nil
}
