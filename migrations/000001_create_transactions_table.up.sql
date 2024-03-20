CREATE TABLE IF NOT EXISTS transactions (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    amount_btc float NOT NULL,
    price_per_btc integer NOT NULL,
    transaction_type integer NOT NULL,
    note text,
    version integer NOT NULL DEFAULT 1
);