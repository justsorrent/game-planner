-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE if not exists users(
    id uuid PRIMARY KEY references user_credentials(id),
    display_name varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE users;
-- +goose StatementEnd
