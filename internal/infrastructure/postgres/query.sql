-- name: CreateExpense :one
INSERT INTO expenses (
    amount,
    category,
    description,
    budget
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING
    id,
    amount,
    category,
    description,
    created_at,
    budget;


-- name: GetExpenseByID :one
SELECT
    id,
    amount,
    category,
    description,
    created_at,
    budget
FROM expenses
WHERE id = $1;


-- name: GetAllExpenses :many
SELECT
    id,
    amount,
    category,
    description,
    created_at,
    budget
FROM expenses
ORDER BY created_at DESC;


-- name: GetAllExpensesByCategory :many
SELECT
    id,
    amount,
    category,
    description,
    created_at,
    budget
FROM expenses
WHERE category = $1
ORDER BY created_at DESC;


-- name: GetAllExpensesByDate :many
SELECT
    id,
    amount,
    category,
    description,
    created_at,
    budget
FROM expenses
WHERE created_at::date = $1
ORDER BY created_at DESC;


-- name: GetAllExpensesByAmount :many
SELECT
    id,
    amount,
    category,
    description,
    created_at,
    budget
FROM expenses
WHERE amount = $1
ORDER BY created_at DESC;


-- name: DeleteExpense :exec
DELETE FROM expenses
WHERE id = $1;