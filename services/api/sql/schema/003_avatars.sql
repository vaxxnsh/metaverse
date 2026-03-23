-- +goose Up

CREATE TABLE avatars (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image_url TEXT NOT NULL,
    name      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE users
    ADD COLUMN avatar_id UUID REFERENCES avatars(id) ON DELETE SET NULL;

ALTER TABLE admins
    ADD COLUMN avatar_id UUID REFERENCES avatars(id) ON DELETE SET NULL;

-- +goose Down

ALTER TABLE users DROP COLUMN avatar_id;

DROP TABLE avatars;
