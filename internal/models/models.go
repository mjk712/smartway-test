package models

type Ticket struct {
	TicketId         int    `json:"ticketId" db:"ticket_id"`
	DeparturePoint   string `json:"departurePoint" db:"departure_point"`
	DestinationPoint string `json:"destinationPoint" db:"destination_point"`
	OrderNumber      string `json:"orderNumber" db:"order_number"`
	ServiceProvider  string `json:"serviceProvider" db:"service_provider"`
	DepartureDate    string `json:"departureDate" db:"departure_date"`
	ArrivalDate      string `json:"arrivalDate" db:"arrival_date"`
	PassengerId      int    `json:"passengerId,omitempty" db:"ticket.passenger_id" `
	CreatedAt        string `json:"createdAt" db:"created_at"`
}

type Passenger struct {
	PassengerId int    `json:"passengerId" db:"passenger_id" `
	LastName    string `json:"lastName" db:"last_name"`
	FirstName   string `json:"firstName" db:"first_name"`
	MiddleName  string `json:"middleName" db:"middle_name"`
}

type Document struct {
	DocumentId     int    `json:"documentId" db:"document_id" `
	PassengerId    int    `json:"passengerId" db:"document.passenger_id" `
	DocumentType   string `json:"documentType" db:"document_type"`
	DocumentNumber string `json:"documentNumber" db:"document_number"`
}
