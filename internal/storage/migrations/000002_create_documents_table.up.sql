CREATE TABLE IF NOT EXISTS document (
                                        document_id BIGSERIAL PRIMARY KEY,
                                        passenger_id INTEGER NOT NULL ,
                                        document_type VARCHAR(50) NOT NULL,
                                        document_number VARCHAR(50) NOT NULL,
                                        FOREIGN KEY (passenger_id) REFERENCES passenger(passenger_id)
);