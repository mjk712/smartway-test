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
