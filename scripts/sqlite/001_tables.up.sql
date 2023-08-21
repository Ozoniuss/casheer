--todo: add indexes

-- Note that INTEGER PRIMARY KEY in sqlite3 is an alias to ROWID.
-- https://www.sqlite.org/autoinc.html

-- Todo: indexes
-- Todo: maybe index just entry id to not also index expense id?

CREATE TABLE IF NOT EXISTS entries(
    id INTEGER PRIMARY KEY,
    month INTEGER NOT NULL, -- smallint?
    year INTEGER NOT NULL,
    category TEXT NOT NULL COLLATE NOCASE,
    subcategory TEXT NOT NULL COLLATE NOCASE,
    -- Even though these values are not necessarily integers, I only care about
    -- two digits of precision, so I will use "bani" or "cents" as the lowest
    -- unit of precis   ion. I'm still debating between dollar, euro and leu, but
    -- will likely go with euro.
    expected_total INTEGER NOT NULL,
    running_total INTEGER NOT NULL,
    recurring BOOLEAN NOT NULL DEFAULT FALSE, -- this may not even be needed
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- do I really need timestamps tho?
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME, -- safety net
    CONSTRAINT unique_logical_entry UNIQUE (month, year, category, subcategory)
);

CREATE TABLE IF NOT EXISTS expenses(
    id INTEGER PRIMARY KEY,
    entry_id INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    currency TEXT NOT NULL, -- iso 4217 currency code
    exponent INTEGER NOT NULL,
    name TEXT NOT NULL COLLATE NOCASE,
    description TEXT COLLATE NOCASE,
    payment_method TEXT COLLATE NOCASE,
    paid_at DATETIME, -- sometimes I might find it useful to know when the expense was paid at
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    CONSTRAINT entry_key FOREIGN KEY(entry_id) REFERENCES entries(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS debts(
    id INTEGER PRIMARY KEY,
    person TEXT NOT NULL COLLATE NOCASE,
    amount INTEGER NOT NULL,
    details TEXT COLLATE NOCASE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME -- safety net
);
