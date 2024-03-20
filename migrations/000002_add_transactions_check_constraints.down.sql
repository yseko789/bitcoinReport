ALTER TABLE transactions DROP CONSTRAINT IF EXISTS btc_amount_check;

ALTER TABLE transactions DROP CONSTRAINT IF EXISTS price_per_btc_check_check;

ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transaction_type_check;
