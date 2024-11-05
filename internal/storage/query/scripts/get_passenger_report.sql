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
