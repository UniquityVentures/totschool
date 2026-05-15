-- +goose Up
INSERT INTO roles (name, created_at, updated_at)
SELECT 'totschool_student', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = 'totschool_student' AND deleted_at IS NULL);

INSERT INTO roles (name, created_at, updated_at)
SELECT 'totschool_admin', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = 'totschool_admin' AND deleted_at IS NULL);

-- +goose Down
DELETE FROM roles WHERE name IN ('totschool_student', 'totschool_admin');
