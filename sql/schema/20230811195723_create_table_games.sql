-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE if not exists games (
    id uuid PRIMARY KEY,
    name varchar(255) NOT NULL,
    description text,
    url varchar(255),
    starting_at timestamp,
    ending_at timestamp,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE games;
-- +goose StatementEnd
