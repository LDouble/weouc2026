CREATE TABLE IF NOT EXISTS iam_users (
    id TEXT PRIMARY KEY,
    open_id TEXT NOT NULL UNIQUE,
    nickname TEXT NOT NULL,
    avatar_url TEXT NOT NULL DEFAULT '',
    roles JSONB NOT NULL DEFAULT '[]'::jsonb,
    permissions JSONB NOT NULL DEFAULT '[]'::jsonb,
    student_profile JSONB,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
