-- +goose Up
-- +goose StatementBegin
ALTER TABLE project ADD COLUMN failure_patterns TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project DROP COLUMN failure_patterns;
-- +goose StatementEnd
