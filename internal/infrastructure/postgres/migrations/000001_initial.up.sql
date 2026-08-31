CREATE TABLE expenses (
    id BIGSERIAL PRIMARY KEY,
    amount NUMERIC(10, 2) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    budget NUMERIC(10, 2) NOT NULL
);