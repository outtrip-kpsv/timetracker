CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS person (
                                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                      surname VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255),
    address VARCHAR(255) NOT NULL,
    passport_number VARCHAR(255) NOT NULL UNIQUE
    );

CREATE INDEX idx_person_surname ON person(surname);
CREATE INDEX idx_person_name ON person(name);
CREATE UNIQUE INDEX idx_person_passport_number ON person(passport_number);