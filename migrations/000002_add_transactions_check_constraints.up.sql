ALTER TABLE transactions ADD CONSTRAINT btc_amount_check CHECK (amount_btc > 0);

ALTER TABLE transactions ADD CONSTRAINT price_per_btc_check CHECK (price_per_btc >= 0);

ALTER TABLE transactions ADD CONSTRAINT transaction_type_check CHECK (transaction_type=1 OR transaction_type=2)

