CREATE TABLE IF NOT EXISTS tasks (
                                     id SERIAL PRIMARY KEY,
                                     title TEXT NOT NULL,
                                     description TEXT,
                                     is_done BOOLEAN DEFAULT false,
                                     created_at TIMESTAMP NOT NULL,
                                     updated_at TIMESTAMP NOT NULL
);
