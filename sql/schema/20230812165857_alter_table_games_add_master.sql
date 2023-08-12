-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE games add column gm_id uuid references users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE games drop column gm_id;
-- +goose StatementEnd
