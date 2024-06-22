package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Partner2 struct {
	BaseURL string
}

type Partner2ReservationRequest struct {
	Lugares      []string `json:"lugares"`
	TipoIngresso string   `json:"tipo_ingresso"`
	Email        string   `json:"email"`
}

type Partner2ReservationResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Lugar        string `json:"lugar"`
	TipoIngresso string `json:"tipo_ingresso"`
	Estado       string `json:"Estado"`
	EventID      string `json:"evento_id"`
}

func (p *Partner2) MakeReservation(req *ReservationRequest) ([]ReservationResponse, error) {
	// Instanciate partnerRequest
	partnerRequest := Partner2ReservationRequest{
		Lugares:      req.Spots,
		TipoIngresso: req.TicketType,
		Email:        req.Email,
	}

	// Convert body to json
	body, err := json.Marshal(partnerRequest)
	if err != nil {
		return nil, err
	}

	// Create http call
	url := fmt.Sprintf("%s/eventos/%s/reservar", p.BaseURL, req.EventID)
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
	var partnerResponse []Partner2ReservationResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&partnerResponse); err != nil {
		return nil, err
	}

	// Convert Partner2ReservationResponse to ReservationResponse
	responses := make([]ReservationResponse, len(partnerResponse))
	for i, r := range partnerResponse {
		responses[i] = ReservationResponse{
			ID:     r.ID,
			Spot:   r.Lugar,
			Status: r.Estado,
		}
	}

	return responses, nil
}
