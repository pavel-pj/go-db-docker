CREATE TABLE payments(
    id SERIAL PRIMARY KEY,
    invoice_id TEXT NOT NULL,
    amount_cents INTEGER NOT NULL,
    status TEXT ,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);