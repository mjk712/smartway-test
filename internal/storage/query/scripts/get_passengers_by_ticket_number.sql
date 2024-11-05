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