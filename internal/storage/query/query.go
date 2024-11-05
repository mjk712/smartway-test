package query

var GetTickets = `SELECT
    ticket_id,
    departure_point,
    destination_point,
    order_number,
    service_provider,
    departure_date,
    arrival_date,
    passenger_id AS "ticket.passenger_id",
    created_at
    FROM flight_ticket;`

var GetPassengersByTicketNumber = `
SELECT 
    p.passenger_id,
    p.last_name,
    p.first_name,
    p.middle_name
FROM 
    flight_ticket t
JOIN 
    passenger p ON t.passenger_id = p.passenger_id
WHERE 
    t.order_number = $1;
`
var GetDocumentsByPassengerId = `
SELECT 
    document_id,document_type,document_number, passenger_id AS "document.passenger_id"
FROM document
WHERE passenger_id = $1;`

var GetTicketFullInfo = `
SELECT 
    ft.ticket_id,
    ft.departure_point,
    ft.destination_point,
    ft.order_number,
    ft.service_provider,
    ft.departure_date,
    ft.arrival_date,
    ft.created_at,
    p.passenger_id,
    p.last_name,
    p.first_name,
    p.middle_name,
    d.document_id,
    d.passenger_id AS "document.passenger_id",
    d.document_type,
    d.document_number
FROM 
    flight_ticket ft
JOIN 
    passenger p ON ft.passenger_id = p.passenger_id
LEFT JOIN 
    document d ON p.passenger_id = d.passenger_id
WHERE 
    ft.order_number = $1;

`
var DeleteTicket = `
DELETE FROM flight_ticket WHERE ticket_id = $1;
`
var DeletePassenger = `
DELETE FROM passenger WHERE passenger_id = $1;
`
var DeleteDocument = `
DELETE FROM document WHERE document_id = $1;
`

var GetPassengerReport = `
SELECT 
    ft.ticket_id,
    ft.departure_point,
    ft.destination_point,
    ft.order_number,
    ft.service_provider,
    ft.departure_date,
    ft.arrival_date,
    ft.created_at,
    CASE
        WHEN ft.created_at < $1 AND ft.departure_date BETWEEN $1 AND $2 THEN 'Ordered earlier, flown in period'
        WHEN ft.created_at BETWEEN $1 AND $2 AND ft.departure_date > $2 THEN 'Ordered in period, not flown'
        WHEN ft.created_at BETWEEN $1 AND $2 AND ft.departure_date BETWEEN $1 AND $2 THEN 'Ordered and flown in period'
    END AS flight_status
FROM 
    flight_ticket ft
WHERE 
    ft.passenger_id = $3
  AND (
        (ft.created_at < $1 AND ft.departure_date BETWEEN $1 AND $2) OR
        (ft.created_at BETWEEN $1 AND $2 AND ft.departure_date > $2) OR
        (ft.created_at BETWEEN $1 AND $2 AND ft.departure_date BETWEEN $1 AND $2)
    )
ORDER BY 
    ft.departure_date;

`
