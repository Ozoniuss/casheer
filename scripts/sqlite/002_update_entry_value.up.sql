BEGIN TRANSACTION;

ALTER TABLE entries
RENAME COLUMN expected_total TO amount;

-- Adding defaults here because I wanted default anyway so it's shorter to 
-- write to the db directly. 
ALTER TABLE entries
ADD COLUMN exponent INTEGER NOT NULL DEFAULT -2;

ALTER TABLE entries
ADD COLUMN currency TEXT COLLATE NOCASE NOT NULL DEFAULT 'RON';

COMMIT;