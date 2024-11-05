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

-- Индекс для поиска по номеру заказа в таблице flight_ticket
CREATE INDEX IF NOT EXISTS idx_flight_ticket_order_number ON flight_ticket (order_number);
-- Индекс для поиска по дате заказа (created_at) в таблице flight_ticket
CREATE INDEX IF NOT EXISTS idx_flight_ticket_created_at ON flight_ticket (created_at);
-- Индекс для поиска по дате вылета (departure_date) в таблице flight_ticket
CREATE INDEX IF NOT EXISTS idx_flight_ticket_departure_date ON flight_ticket (departure_date);