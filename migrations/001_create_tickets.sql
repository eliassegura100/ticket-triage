CREATE TABLE IF NOT EXISTS tickets (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL CHECK (btrim(title) <> ''), --Ensure title is not empty or whitespace
    description TEXT NOT NULL DEFAULT '',
    reporter TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'triaged', 'resolved')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);