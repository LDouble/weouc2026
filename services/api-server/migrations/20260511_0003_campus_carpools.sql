CREATE TABLE IF NOT EXISTS campus_carpools (
    id TEXT PRIMARY KEY,
    category TEXT NOT NULL DEFAULT '',
    route_from TEXT NOT NULL DEFAULT '',
    route_to TEXT NOT NULL DEFAULT '',
    travel_at TIMESTAMPTZ NOT NULL,
    type_label TEXT NOT NULL DEFAULT '',
    seats_text TEXT NOT NULL DEFAULT '',
    price_text TEXT NOT NULL DEFAULT '',
    note TEXT NOT NULL DEFAULT '',
    tags JSONB NOT NULL DEFAULT '[]'::jsonb,
    contact TEXT NOT NULL DEFAULT '',
    review_status TEXT NOT NULL DEFAULT 'published',
    publisher_user_id TEXT NOT NULL,
    publisher TEXT NOT NULL,
    publisher_initial TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);
