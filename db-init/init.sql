
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title text,
    description text,
    status text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);