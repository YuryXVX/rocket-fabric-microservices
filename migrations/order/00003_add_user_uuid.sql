-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD COLUMN user_uuid UUID;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders DROP COLUMN IF EXISTS user_uuid;
-- +goose StatementEnd