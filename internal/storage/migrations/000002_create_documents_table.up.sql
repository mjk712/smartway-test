CREATE TABLE IF NOT EXISTS document (
                                         id BIGSERIAL PRIMARY KEY,
                                         passenger_id INTEGER NOT NULL ,
                                         last_name VARCHAR(50) NOT NULL,
                                         middle_name VARCHAR(50) NOT NULL,
                                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                         FOREIGN KEY (passenger_id) REFERENCES passenger(id)
);