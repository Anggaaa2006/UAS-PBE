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
    GetByEmail(ctx context.Context, email string) (*model.User, error) // alias
    Create(ctx context.Context, user model.User) error
    List(ctx context.Context) ([]model.User, error)
    GetByID(ctx context.Context, id string) (*model.User, error)
    Update(ctx context.Context, user model.User) error
    Delete(ctx context.Context, id string) error
    UpdateRole(ctx context.Context, id string, role string) error

    
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
    FindByEmail
    Ambil user berdasarkan email (primary method)
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
    GetByEmail
    Alias untuk FindByEmail agar service bisa pakai
*/
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    return r.FindByEmail(ctx, email)
}

/*
    Create user (register)
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

/*
    List semua user
*/
func (r *userRepo) List(ctx context.Context) ([]model.User, error) {

    q := `
        SELECT id, name, email, role
        FROM users
        ORDER BY name
    `

    rows, err := r.db.QueryContext(ctx, q)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []model.User

    for rows.Next() {
        var u model.User

        err := rows.Scan(
            &u.ID,
            &u.Name,
            &u.Email,
            &u.Role,
        )
        if err != nil {
            return nil, err
        }

        list = append(list, u)
    }

    return list, nil
}

/* existing struct and constructor remain unchanged */

/*
    GetByID
    Mengambil data user berdasarkan ID
*/
func (r *userRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, name, email, role, created_at
		FROM users
		WHERE id = $1
	`

	var u model.User
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt)

    if 

	if err != nil {
		return nil, err
	}

	return &user, nil
}
/*
    Update user (name, email, role optionally, password must be hashed before)
*/
func (r *userRepo) Update(ctx context.Context, user model.User) error {
    q := `UPDATE users SET name = $1, email = $2, password = $3, role = $4, created_at = COALESCE(created_at, NOW()) WHERE id = $5`
    _, err := r.db.ExecContext(ctx, q, user.Name, user.Email, user.Password, user.Role, user.ID)
    return err
}

/*
    Delete user
*/
func (r *userRepo) Delete(ctx context.Context, id string) error {
    q := `DELETE FROM users WHERE id = $1`
    _, err := r.db.ExecContext(ctx, q, id)
    return err
}

/*
    UpdateRole (patch role)
*/
func (r *userRepo) UpdateRole(ctx context.Context, id string, role string) error {
    q := `UPDATE users SET role = $1 WHERE id = $2`
    _, err := r.db.ExecContext(ctx, q, role, id)
    return err
}