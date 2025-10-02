-- +goose Up
ALTER TABLE posts 
ALTER COLUMN published_at DROP NOT NULL;

-- +goose Down
ALTER TABLE posts 
ALTER COLUMN published_at SET NOT NULL;
