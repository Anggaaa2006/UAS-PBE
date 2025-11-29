package postgres

import (
	"context"
	"database/sql"
	"uas_pbe/model"
)

type RoleRepo interface {
	GetByID(ctx context.Context, id string) (*model.Role, error)
}

type roleRepo struct{ db *sql.DB }

func NewRoleRepo(db *sql.DB) RoleRepo {
	return &roleRepo{db}
}

func (r *roleRepo) GetByID(ctx context.Context, id string) (*model.Role, error) {
	query := `SELECT id, name FROM roles WHERE id=$1`
	row := r.db.QueryRowContext(ctx, query, id)

	var role model.Role
	if err := row.Scan(&role.ID, &role.Name); err != nil {
		return nil, err
	}

	return &role, nil
}
