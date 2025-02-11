package repository

type UserRepo interface {
	Add(ctx context.Context, u *models.User)
}

type PostgresUserRepo struct {
	DB *sql.DB
}

func NewPostgreUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{
		DB: db,
	}
}

func (r *PostgresUserRepo) Add(ctx context.Context, u *models.User, token string) error {
	query := `
	    INSERT INTO users(username, password)
	    VALUES ($1, $2)
	    RETURNING id, created_at`

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	
	args := []any{u.Username, u.Password}

	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		tx.RollBack()
		return err
	}

	query = `
	     INSERT INTO auth_token(user_id, token)
	     VALUES ($1, $2)`
	
	err = 
}
