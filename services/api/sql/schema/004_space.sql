-- +goose Up

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE spaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER,
    thumbnail TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE elements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    image_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE maps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE space_elements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    space_id UUID NOT NULL,
    element_id UUID NOT NULL,
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_space
        FOREIGN KEY(space_id)
        REFERENCES spaces(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_element
        FOREIGN KEY(element_id)
        REFERENCES elements(id)
        ON DELETE CASCADE
);

CREATE TABLE map_elements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    map_id UUID NOT NULL,
    element_id UUID,
    x INTEGER,
    y INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_map
        FOREIGN KEY(map_id)
        REFERENCES maps(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_element_optional
        FOREIGN KEY(element_id)
        REFERENCES elements(id)
);

CREATE INDEX idx_space_elements_space_id ON space_elements(space_id);
CREATE INDEX idx_space_elements_element_id ON space_elements(element_id);

CREATE INDEX idx_map_elements_map_id ON map_elements(map_id);
CREATE INDEX idx_map_elements_element_id ON map_elements(element_id);

CREATE UNIQUE INDEX unique_space_position
ON space_elements(space_id, x, y);

CREATE INDEX idx_space_elements_space_xy
ON space_elements(space_id, x, y);

-- +goose Down

DROP TABLE IF EXISTS map_elements;
DROP TABLE IF EXISTS space_elements;
DROP TABLE IF EXISTS maps;
DROP TABLE IF EXISTS elements;
DROP TABLE IF EXISTS spaces;