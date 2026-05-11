ALTER TABLE campus_markets
    ADD COLUMN IF NOT EXISTS review_status TEXT NOT NULL DEFAULT 'published';

ALTER TABLE campus_errands
    ADD COLUMN IF NOT EXISTS review_status TEXT NOT NULL DEFAULT 'published';

ALTER TABLE campus_resources
    ADD COLUMN IF NOT EXISTS review_status TEXT NOT NULL DEFAULT 'published';

ALTER TABLE campus_lost_founds
    ADD COLUMN IF NOT EXISTS review_status TEXT NOT NULL DEFAULT 'published';
