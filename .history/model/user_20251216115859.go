package model

/*
    Model User untuk login/auth
*/
type User struct {
	ID       string `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"`  // tidak boleh ditampilkan kembali
	Role     string `db:"role" json:"role"`   // WAJIB, dipakai untuk JWT & middleware
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
