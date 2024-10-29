CREATE TABLE IF NOT EXISTS passenger (
                                        id BIGSERIAL PRIMARY KEY ,
                                        last_name VARCHAR(50) NOT NULL,
                                        middle_name VARCHAR(50) NOT NULL,
                                        first_name VARCHAR(50),
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);