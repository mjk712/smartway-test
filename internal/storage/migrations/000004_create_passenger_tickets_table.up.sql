CREATE TABLE IF NOT EXISTS passenger_ticket (
                                        id BIGSERIAL PRIMARY KEY,
                                        passenger_id INTEGER NOT NULL ,
                                        ticket_id INTEGER NOT NULL ,
                                        FOREIGN KEY (passenger_id) REFERENCES passenger(passenger_id),
                                        FOREIGN KEY (ticket_id) REFERENCES flight_ticket(ticket_id)
);