-- +goose Up
-- +goose StatementBegin
INSERT INTO apps (id, name)
VALUES (1, 'test')
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
