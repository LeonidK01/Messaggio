-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS message (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	created_by uuid NOT NULL,
	updated_by uuid NOT NULL,
	deleted_by uuid,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at timestamp,

    from uuid NOT NULL,
    to uuid NOT NULL,
	text text NOT NULL,
	sended boolean NOT NULL DEFAULT false,

	CONSTRAINT mesage_pkey PRIMARY KEY (id)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS message;
