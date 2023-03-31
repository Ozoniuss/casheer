--todo: add indexes

CREATE TABLE IF NOT EXISTS entries(
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    month SMALLINT NOT NULL,
    year SMALLINT NOT NULL,
    category TEXT NOT NULL,
    subcategory TEXT NOT NULL,
    expected_total REAL NOT NULL,
    running_total REAL NOT NULL,
    recurring BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT unique_logical_entry UNIQUE (month, year, category, subcategory)
);

CREATE TABLE IF NOT EXISTS expenses(
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL,
    value REAL NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    payment_method TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT entry_key FOREIGN KEY(entry_id) REFERENCES entries(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS debts(
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    person TEXT NOT NULL,
    amount REAL NOT NULL,
    details TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);
