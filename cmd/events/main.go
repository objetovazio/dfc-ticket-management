package main

import (
	"database/sql"
	"net/http"

	httpHandler "github.com/devfullcycle/imersao18/golang/internal/events/domain/infra/http"
	"github.com/devfullcycle/imersao18/golang/internal/events/domain/infra/service"
	"github.com/devfullcycle/imersao18/golang/internal/events/domain/infra/service/repository"
	"github.com/devfullcycle/imersao18/golang/internal/events/usecase"
)

func main() {
	db, err := sql.Open("mysql", "test_user:password@tcp(127.0.0.1:3306)/test_database)")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventRepo, err := repository.NewMysqlEventRepository(db)
	if err != nil {
		panic(err)
	}

	partnerBaseURLs := map[int]string{
		1: "http://localhost:3333",
		3: "http://localhost:3334",
	}

	partnerFactory := service.NewPartnerFactory(partnerBaseURLs)

	listEventsUseCase := usecase.NewListEventsUseCase(eventRepo)
	getEventsUseCase := usecase.NewGetEventUseCase(eventRepo)
	listSpotsUseCase := usecase.NewListSpotsUseCase(eventRepo)
	buyTicketUseCase := usecase.NewBuyTicketsUseCase(eventRepo, partnerFactory)

	eventsHandler := httpHandler.NewEventHandler(
		listEventsUseCase,
		listSpotsUseCase,
		getEventsUseCase,
		buyTicketUseCase,
	)

	r := http.NewServeMux()
	r.HandleFunc("/events", eventsHandler.ListEvents)
	r.HandleFunc("/events/{eventID}", eventsHandler.GetEvents)
	r.HandleFunc("/events/{eventID}", eventsHandler.GetEvents)
	r.HandleFunc("/events/{eventID}/spots", eventsHandler.ListSpots)
	r.HandleFunc("POST /checkout", eventsHandler.BuyTickets)
	http.ListenAndServe(":8080", r)
}
