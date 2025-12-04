package postgres

import (
    "context"
    "database/sql"
    "errors"

    "uas_pbe/model"
)

/*
    Interface UserRepo
*/
type UserRepo interface {
    FindByEmail(ctx context.Context, email string) (*model.User, error)
    Create(ctx context.Context, user model.User) error
}

/*
    Struct implementasi UserRepo
*/
type userRepo struct {
    db *sql.DB
}

/*
    Constructor
*/
func NewUserRepo(db *sql.DB) UserRepo {
    return &userRepo{db: db}
}

/*
    Ambil user berdasarkan email
*/
func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {

    q := `
        SELECT id, name, email, password, role
        FROM users
        WHERE email = $1
        LIMIT 1
    `

    var user model.User

    err := r.db.QueryRowContext(ctx, q, email).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Password,
        &user.Role,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }

    return &user, nil
}

/*
    Create user (opsional, untuk register)
*/
func (r *userRepo) Create(ctx context.Context, user model.User) error {

    q := `
        INSERT INTO users (id, name, email, password, role)
        VALUES ($1, $2, $3, $4, $5)
    `

    _, err := r.db.ExecContext(ctx, q,
        user.ID,
        user.Name,
        user.Email,
        user.Password,
        user.Role,
    )

    return err
}
