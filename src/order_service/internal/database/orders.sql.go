// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: orders.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (id, created_at, modified_at, user_id, product_id) 
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, modified_at, user_id, product_id
`

type CreateOrderParams struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
	UserID     uuid.UUID
	ProductID  uuid.UUID
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, createOrder,
		arg.ID,
		arg.CreatedAt,
		arg.ModifiedAt,
		arg.UserID,
		arg.ProductID,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.ModifiedAt,
		&i.UserID,
		&i.ProductID,
	)
	return i, err
}

const deleteOrder = `-- name: DeleteOrder :exec
DELETE FROM orders WHERE id=$1 AND user_id=$2
`

type DeleteOrderParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteOrder(ctx context.Context, arg DeleteOrderParams) error {
	_, err := q.db.ExecContext(ctx, deleteOrder, arg.ID, arg.UserID)
	return err
}

const getOrderById = `-- name: GetOrderById :one
SELECT id, created_at, modified_at, product_id FROM orders WHERE id=$1 AND user_id=$2
`

type GetOrderByIdParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

type GetOrderByIdRow struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
	ProductID  uuid.UUID
}

func (q *Queries) GetOrderById(ctx context.Context, arg GetOrderByIdParams) (GetOrderByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getOrderById, arg.ID, arg.UserID)
	var i GetOrderByIdRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.ModifiedAt,
		&i.ProductID,
	)
	return i, err
}

const getOrders = `-- name: GetOrders :many
SELECT id, product_id FROM orders WHERE user_id=$1 LIMIT $2 OFFSET $3
`

type GetOrdersParams struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

type GetOrdersRow struct {
	ID        uuid.UUID
	ProductID uuid.UUID
}

func (q *Queries) GetOrders(ctx context.Context, arg GetOrdersParams) ([]GetOrdersRow, error) {
	rows, err := q.db.QueryContext(ctx, getOrders, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrdersRow
	for rows.Next() {
		var i GetOrdersRow
		if err := rows.Scan(&i.ID, &i.ProductID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
