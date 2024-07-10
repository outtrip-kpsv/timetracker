CREATE TYPE status AS ENUM ('new', 'work', 'pause', 'complete');

CREATE TABLE tasks (
                       IdTask UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       IdPerson UUID NOT NULL,
                       task_name TEXT NOT NULL,
                       task_status status NOT NULL,
                       FOREIGN KEY (IdPerson) REFERENCES person (id)
);