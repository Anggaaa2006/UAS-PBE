/*
	Create
	Membuat user baru (dipakai saat register)
*/
func (r *userRepo) Create(ctx context.Context, u model.User) error {
	query := `
		INSERT INTO users (id, name, email, password, role_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Name, u.Email, u.Password, u.RoleID)
	return err
}
