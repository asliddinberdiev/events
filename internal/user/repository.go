package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/asliddinberdiev/events/internal/common"
	"github.com/asliddinberdiev/events/pkg/db"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, req RegisterPayload) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByTelegramID(ctx context.Context, telegram_id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetAll(ctx context.Context, params common.SearchRequest) (*common.AllResponse, error)
	Update(ctx context.Context, req UpdateUser) error
	Delete(ctx context.Context, user_id uuid.UUID) error
}

type PostgresUserRepository struct{}

func NewPostgresUserRepository() *PostgresUserRepository {
	return &PostgresUserRepository{}
}

func (r *PostgresUserRepository) Create(ctx context.Context, req RegisterPayload) (*User, error) {
	query := `
		INSERT INTO users (telegram_id, username,	password)
		VALUES ($1, $2,	$3)
		RETURNING id, telegram_id, username, password
	`

	var response User
	if err := db.DB.QueryRowContext(ctx, query, req.TelegramID, req.Username, req.Password).Scan(&response.ID, &response.TelegramID, &response.Username, &response.Password); err != nil {
		return nil, errors.Wrap(err, "Error: create user")
	}

	return &response, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, telegram_id, username, password
		FROM users
		WHERE id = $1
	`

	var response User
	if err := db.DB.GetContext(ctx, &response, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "Error: get user by id")
	}

	return &response, nil
}

func (r *PostgresUserRepository) GetByTelegramID(ctx context.Context, telegram_id string) (*User, error) {
	query := `
		SELECT id, telegram_id, username, password
		FROM users
		WHERE telegram_id = $1
	`

	var response User
	if err := db.DB.GetContext(ctx, &response, query, telegram_id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "Error: get user by telegram_id")
	}

	return &response, nil
}

func (r *PostgresUserRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, telegram_id, username, password
		FROM users
		WHERE username = $1
	`

	var response User
	if err := db.DB.GetContext(ctx, &response, query, username); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "Error: get user by username")
	}

	return &response, nil
}

func (r *PostgresUserRepository) GetAll(ctx context.Context, params common.SearchRequest) (*common.AllResponse, error) {
	allusers := &common.AllResponse{
		Total: 0,
	}

	offset := (params.Page - 1) * params.Limit
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, offset)

	filter := ""
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(" AND username ILIKE '%s'", str)
	}

	totalQuery := `
    SELECT COUNT(*) 
    FROM users ` + filter
	var totalCount int
	if err := db.DB.GetContext(ctx, &totalCount, totalQuery); err != nil {
		return nil, errors.Wrap(err, "Error: get total count of users")
	}
	allusers.Total = uint32(totalCount)

	query := `
		SELECT id, telegram_id, username
		FROM users ` + filter + " " + limit

	var list []*User
	if err := db.DB.SelectContext(ctx, &list, query); err != nil {
		return nil, errors.Wrap(err, "Error: get all users")
	}

	allusers.List = make([]any, len(list))
	for i, user := range list {
		allusers.List[i] = user
	}

	return allusers, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, req UpdateUser) error {
	query := `
		UPDATE users
		SET 
			username = $1
		WHERE id = $2
	`

	_, err := db.DB.ExecContext(ctx, query, req.Username, req.ID)
	if err != nil {
		return errors.Wrap(err, "Error: user update")
	}

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, user_id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := db.DB.ExecContext(ctx, query, user_id)
	return err
}
