BEGIN TRANSACTION;

ALTER TABLE entries
RENAME COLUMN amount TO expected_total;

ALTER TABLE entries
DROP COLUMN exponent;

ALTER TABLE entries
DROP COLUMN currency;

COMMIT;
