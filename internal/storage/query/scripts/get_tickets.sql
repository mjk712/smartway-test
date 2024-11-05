SELECT
    ticket_id,
    departure_point,
    destination_point,
    order_number,
    service_provider,
    departure_date,
    arrival_date,
    passenger_id AS "ticket.passenger_id",
    created_at
FROM flight_ticket;