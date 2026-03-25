-- +goose Up

CREATE TABLE elements (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    width      INTEGER NOT NULL,
    height     INTEGER NOT NULL,
    image_url  TEXT NOT NULL,
    static     BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE maps (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL,
    width      INTEGER NOT NULL,
    height     INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE map_elements (
    map_id     UUID NOT NULL REFERENCES maps(id) ON DELETE CASCADE,
    element_id UUID NOT NULL REFERENCES elements(id) ON DELETE CASCADE,
    x          INTEGER NOT NULL,
    y          INTEGER NOT NULL,
    PRIMARY KEY (map_id, x, y)
);

CREATE TABLE spaces (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    width       INTEGER NOT NULL,
    height      INTEGER NOT NULL,
    thumbnail   TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE space_elements (
    space_id   UUID NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
    element_id UUID NOT NULL REFERENCES elements(id) ON DELETE CASCADE,
    x          INTEGER NOT NULL,
    y          INTEGER NOT NULL,
    PRIMARY KEY (space_id, x, y)
);

-- +goose Down

DROP TABLE IF EXISTS space_elements;
DROP TABLE IF EXISTS spaces;
DROP TABLE IF EXISTS map_elements;
DROP TABLE IF EXISTS maps;
DROP TABLE IF EXISTS elements;