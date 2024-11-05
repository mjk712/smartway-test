package requests

type TicketUpdateRequest struct {
	DeparturePoint   *string `json:"departurePoint" `
	DestinationPoint *string `json:"destinationPoint" `
	ServiceProvider  *string `json:"serviceProvider" `
	DepartureDate    *string `json:"departureDate" `
	ArrivalDate      *string `json:"arrivalDate" `
	PassengerId      *int    `json:"passengerId" `
	CreatedAt        *string `json:"createdAt" `
}

type UpdatePassengerRequest struct {
	LastName   *string `json:"lastName"`
	FirstName  *string `json:"firstName"`
	MiddleName *string `json:"middleName"`
}

type DocumentUpdateRequest struct {
	PassengerId    *int    `json:"passengerId"  `
	DocumentType   *string `json:"documentType"`
	DocumentNumber *string `json:"documentNumber"`
}
