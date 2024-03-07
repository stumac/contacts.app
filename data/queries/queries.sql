-- name: GetContact :one
SELECT * FROM contacts WHERE id = ?;

-- name: ListContacts :many
SELECT * FROM contacts ORDER BY first_name;

-- name:FindContactByProper

-- name: CreateContact :one
INSERT INTO contacts (
    first_name, phone, email
) VALUES (
    ?, ?, ?
) RETURNING *;

-- name: UpdateContact :one
UPDATE contacts
set first_name = ?, phone = ?, email = ? WHERE id = ? RETURNING *;