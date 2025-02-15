-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user"
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username    TEXT UNIQUE NOT NULL CHECK ( username <> '' ),
    email       TEXT UNIQUE NOT NULL CHECK ( email <> '' ),
    password    TEXT NOT NULL CHECK ( octet_length(password) <> 0 ),
    role        TEXT NOT NULL DEFAULT 'user',
    avatar      TEXT,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    login_date  TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user" CASCADE;
-- +goose StatementEnd
