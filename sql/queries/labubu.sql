-- name: GetAllLabubu :many
SELECT * FROM labubu ORDER BY id;

-- name: CreateLabubu :one
INSERT INTO labubu (text) VALUES ($1) RETURNING *;