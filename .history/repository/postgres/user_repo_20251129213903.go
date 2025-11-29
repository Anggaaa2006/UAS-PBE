package postgres

import (
	"context"
	"database/sql"
	"errors"
	"uas_pbe/model"
)

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context) ([]model.User, error)
}

type userRepo struct{ db *sql.DB }

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{db}
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, name, email, password, role_id FROM users WHERE email=$1`
	row := r.db.QueryRowContext(ctx, query, email)

	var u model.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RoleID); err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return &u, nil
}

func (r *userRepo) List(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, email, role_id FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.User
	for rows.Next() {
		var u model.User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.RoleID)
		list = append(list, u)
	}
	return list, nil
}
