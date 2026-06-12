-- name: SaveURL :exec
INSERT INTO urls (code, target) VALUES ($1, $2);

-- name: GetTargetByCode :one
SELECT target FROM urls WHERE code = $1;
