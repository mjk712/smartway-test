CREATE TABLE IF NOT EXISTS flight_ticket (
                                             ticket_id BIGSERIAL PRIMARY KEY,
                                             departure_point VARCHAR(50) NOT NULL,
                                             destination_point VARCHAR(50) NOT NULL,
                                             order_number INTEGER NOT NULL,
                                             service_provider VARCHAR(50) NOT NULL,
                                             departure_date TIMESTAMP NOT NULL,
                                             arrival_date TIMESTAMP NOT NULL,
                                             passenger_id INTEGER NOT NULL ,
                                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                             FOREIGN KEY (passenger_id) REFERENCES passenger(passenger_id)

);