CREATE TABLE IF NOT EXISTS campus_life_id_sequences (
    name TEXT PRIMARY KEY,
    value BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS campus_markets (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    publisher_user_id TEXT NOT NULL,
    publisher TEXT NOT NULL,
    publisher_initial TEXT NOT NULL,
    image TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL,
    likes INTEGER NOT NULL DEFAULT 0,
    liked_by_user_ids JSONB NOT NULL DEFAULT '{}'::jsonb,
    extra JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE IF NOT EXISTS campus_errands (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT '',
    route_start TEXT NOT NULL DEFAULT '',
    route_end TEXT NOT NULL DEFAULT '',
    deadline TIMESTAMPTZ NOT NULL,
    reward TEXT NOT NULL DEFAULT '',
    contact TEXT NOT NULL DEFAULT '',
    urgent BOOLEAN NOT NULL DEFAULT FALSE,
    images JSONB NOT NULL DEFAULT '[]'::jsonb,
    status TEXT NOT NULL DEFAULT 'published',
    publisher_user_id TEXT NOT NULL,
    publisher TEXT NOT NULL,
    publisher_initial TEXT NOT NULL,
    acceptor_user_id TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS campus_resources (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    publisher_user_id TEXT NOT NULL,
    publisher TEXT NOT NULL,
    publisher_initial TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    extra JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE IF NOT EXISTS campus_lost_founds (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    publisher_user_id TEXT NOT NULL,
    publisher TEXT NOT NULL,
    publisher_initial TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    extra JSONB NOT NULL DEFAULT '{}'::jsonb
);
