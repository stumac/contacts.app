-- +goose Up
-- +goose StatementBegin
CREATE VIRTUAL TABLE contacts USING FTS5 (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    phone TEXT,
    email TEXT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE contacts;
-- +goose StatementEnd
