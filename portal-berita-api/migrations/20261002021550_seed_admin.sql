-- +goose Up
INSERT INTO users (nama, email, password, role) VALUES
('Administrator', 'admin@portal.com', '$2b$10$vkUlEW5NcNQNwIwbEwV2se3pBodT8k6.egzzfV9S3A11/NjOdVQea', 'admin');

-- +goose Down
DELETE FROM users WHERE email = 'admin@portal.com';
