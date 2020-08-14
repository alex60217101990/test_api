package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/alex60217101990/test_api/internal/models"
	"github.com/alex60217101990/test_api/internal/repository"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func (r *Repository) GetUserByCreeds(ctx context.Context, creeds *models.Credentials, associate ...struct{}) (user *models.User, err error) {
	defer func() {
		if user != nil && len(user.Password) > 0 {
			user.AfterFind()
		}
		if err != nil {
			err = errors.WithMessage(err, "get user by creeds")
		}
	}()

	var conn *pgxpool.Conn
	conn, err = r.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	if len(associate) > 0 {
		var rows pgx.Rows
		rows, err = conn.Query(ctx,
			fmt.Sprintf(
				`SELECT u.*, pc.*, pr.* FROM public.users u 
			JOIN public.product_categories pc
				ON u.id = pc.change_by_user
			JOIN public.products pr
				ON u.id = pr.change_by_user
			WHERE (u.username = '%s' OR u.email = '%s') 
				AND u.password = '%s'
				AND u.deleted_at IS NULL
				AND pc.deleted_at IS NULL
				AND pr.deleted_at IS NULL;`,
				creeds.Username, creeds.Email, creeds.Password,
			),
		)
		if err != nil {
			return nil, err
		}

		users := make([]*models.User, 0)
		users, err = repository.ConvertMultipleRowsToUser(ctx, rows)
		if err != nil {
			return nil, err
		}
		if len(users) == 0 {
			return nil, pgx.ErrNoRows
		}

		return users[0], err
	}

	row := conn.QueryRow(ctx,
		`SELECT u.* FROM public.users u 
			WHERE (u.username = $1 OR u.email = $2) 
			AND u.password = $3
			AND u.deleted_at IS NULL;`,
		creeds.Username, creeds.Email, creeds.Password,
	)

	return repository.ConvertSingleRowToUser(ctx, row)
}

func (r *Repository) InsertUser(ctx context.Context, user *models.User) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "insert new user posrgres repo")
		}
	}()

	err = user.BeforeCreate()
	if err != nil {
		return err
	}

	var conn *pgxpool.Conn
	conn, err = r.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx,
		`INSERT INTO users (username, email, password, is_online) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (username, email, password) DO 
		UPDATE SET 
			username = EXCLUDED.username, email = EXCLUDED.email, 
			password = EXCLUDED.password, is_online = EXCLUDED.is_online,
			updated_at = (CURRENT_TIMESTAMP at TIME ZONE 'utc');`,
		user.Username, user.Email, user.Password, user.IsOnline,
	)

	return err
}

func (r *Repository) UpdateUser(ctx context.Context, user *models.User) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "update new user posrgres repo")
		}
	}()

	err = user.BeforeCreate()
	if err != nil {
		return err
	}

	var conn *pgxpool.Conn
	conn, err = r.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := fmt.Sprintf(`
	UPDATE users SET
		username = COALESCE('%s', username),
		email = COALESCE('%s', email),
		password = COALESCE('%s', password),
		is_online = COALESCE(%t, is_online),
		updated_at = (CURRENT_TIMESTAMP at TIME ZONE 'utc')
	WHERE public_id = '%s'
	AND ('%s' IS NOT NULL AND '%s' IS DISTINCT FROM username AND (('%s' = '') IS FALSE) 
		OR '%s' IS NOT NULL AND '%s' IS DISTINCT FROM email AND (('%s' = '') IS FALSE)
		OR '%s' IS NOT NULL AND '%s' IS DISTINCT FROM password AND (('%s' = '') IS FALSE)
		OR %t IS NOT NULL AND %t IS DISTINCT FROM is_online
	) 
	AND deleted_at IS NULL;`,
		user.Username, user.Email, user.Password, user.IsOnline, user.PublicID.String(),
		user.Username, user.Username, user.Username,
		user.Email, user.Email, user.Email,
		user.Password, user.Password, user.Password,
		user.IsOnline, user.IsOnline,
	)

	_, err = conn.Exec(ctx, query)

	return err
}

func (r *Repository) DeleteSoft(ctx context.Context, publicID string) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "soft delete user posrgres repo")
		}
	}()

	var conn *pgxpool.Conn
	conn, err = r.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx,
		`UPDATE users SET
		deleted_at = (CURRENT_TIMESTAMP at TIME ZONE 'utc')
	WHERE public_id = $1;`, publicID,
	)

	return err
}

func (r *Repository) DeleteHard(ctx context.Context, publicID string) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "hard delete user posrgres repo")
		}
	}()

	var conn *pgxpool.Conn
	conn, err = r.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx,
		`DELETE FROM users WHERE public_id = $1;`, publicID,
	)

	return err
}
