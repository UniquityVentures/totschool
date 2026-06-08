-- +goose Up
UPDATE clients
SET status = 'archived'::client_status;

-- +goose Down
UPDATE clients
SET status = 'active'::client_status;
