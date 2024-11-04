CREATE TABLE IF NOT EXISTS passenger (
                                        passenger_id BIGSERIAL PRIMARY KEY ,
                                        last_name VARCHAR(50) NOT NULL,
                                        middle_name VARCHAR(50) NOT NULL,
                                        first_name VARCHAR(50)
);