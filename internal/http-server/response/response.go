package response

import "smartway-test/internal/models"

type PassengerWithDocs struct {
	models.Passenger
	Documents []models.Document
}

type FullTicketInfo struct {
	Ticket     models.Ticket
	Passengers []PassengerWithDocs
}

type FlightReport struct {
	models.Ticket
	FlightStatus string `db:"flight_status" json:"flightStatus"`
}
