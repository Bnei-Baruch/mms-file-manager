
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE UNIQUE INDEX file_name_version_uniq_idx ON files (lower(file_name), version);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX file_name_version_uniq_idx;
