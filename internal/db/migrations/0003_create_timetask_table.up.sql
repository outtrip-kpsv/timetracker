CREATE TABLE IF NOT EXISTS timetask (
                                        id SERIAL PRIMARY KEY,
                                        idtask UUID,
                                        start_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        end_time TIMESTAMP,
                                        FOREIGN KEY (idtask) REFERENCES tasks (idtask)
);