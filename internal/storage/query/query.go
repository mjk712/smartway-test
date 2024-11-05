package query

import (
	_ "embed"
)

//go:embed scripts/get_tickets.sql
var GetTickets string

//go:embed scripts/get_passengers_by_ticket_number.sql
var GetPassengersByTicketNumber string

//go:embed scripts/get_documents_by_passenger_id.sql
var GetDocumentsByPassengerId string

//go:embed scripts/get_ticket_full_info.sql
var GetTicketFullInfo string

//go:embed scripts/delete_ticket.sql
var DeleteTicket string

//go:embed scripts/delete_passenger.sql
var DeletePassenger string

//go:embed scripts/delete_document.sql
var DeleteDocument string

//go:embed scripts/get_passenger_report.sql
var GetPassengerReport string
