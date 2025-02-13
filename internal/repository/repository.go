package repository

import (
	"context"
	"database/sql"
	"errors"
	"merch-shop/internal/models"
)

type Repository interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Add(ctx context.Context, u *models.User) error
	GetBalance(ctx context.Context, userID int) (int, error)
	SendCoin(ctx context.Context, senderID, receiverID int, amount int) error
}

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		DB: db,
	}
}

func (r *PostgresRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
	    SELECT id, password_hash, created_at
	    FROM users
	    WHERE username = $1`

	user := &models.User{
		Username: username,
	}

	err := r.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresRepository) Add(ctx context.Context, u *models.User) error {
	query := `
	    INSERT INTO users(username, password_hash)
	    VALUES ($1, $2)
	    RETURNING id, created_at`

	args := []any{u.Username, u.PasswordHash}

	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = tx.QueryRowContext(ctx, query, args...).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `
	    INSERT INTO coins(user_id)
	    VALUES ($1)`

	_, err = tx.ExecContext(ctx, query, u.ID)

	err = tx.Commit()
	return err
}

func (r *PostgresRepository) GetBalance(ctx context.Context, userID int) (int, error) {
	query := `
	    SELECT balance
	    FROM coins
	    WHERE user_id = $1`

	var balance int

	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (r *PostgresRepository) InsertItem(ctx context.Context, userID, itemID int) error {
	query := `
	    INSERT INTO inventory(user_id, item_id)
	    VALUES ($1, $2)`

	args := []any{userID, itemID}

	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	query = `
	    SELECT price
	    FROM item
	    WHERE id = $1`

	var price int

	err = tx.QueryRowContext(ctx, query, itemID).Scan(&price)
	if err != nil {
		return err
	}

	query = `
	    UPDATE coins
	    SET balance = balance - $2
	    WHERE user_id = $1`

	args = []any{userID, price}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (r *PostgresRepository) UpdateItem(ctx context.Context, userID, itemID int) error {
	query := `
	    UPDATE inventory
	    SET quantity = quantity + 1
	    WHERE user_id = $1 AND item_ID = $2`

	args := []any{userID, itemID}

	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	query = `
	    SELECT price
	    FROM item
	    WHERE id = $1`

	var price int

	err = tx.QueryRowContext(ctx, query, itemID).Scan(&price)
	if err != nil {
		return err
	}

	query = `
	    UPDATE coins
	    SET balance = balance - $2
	    WHERE user_id = $1`

	args = []any{userID, price}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (r *PostgresRepository) SendCoin(ctx context.Context, senderID, receiverID int, amount int) error {
	query := `
	     INSERT INTO transaction(sender_id, receiver_id, amount)
	     VALUES ($1, $2, $3)`

	args := []any{senderID, receiverID, amount}

	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	query = `
	    UPDATE coins
	    SET balance = CASE
	    WHEN user_id = $1 THEN balance - $3
	    WHEN user_id = $2 THEN balance + $3
	    ELSE balance
	    END
	    WHERE user_id in ($1, $2)`

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (r *PostgresRepository) GetInventory(ctx context.Context, userID int) ([]*models.InventoryItem, error) {
	query := `
	    SELECT it.type, in.quantity
	    FROM inventory AS in
	    JOIN item AS it ON in.item_id = it.id
	    WHERE user_id = $1`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventory []*models.InventoryItem

	for rows.Next() {
		var item *models.InventoryItem
		err := rows.Scan(
			&item.Type,
			&item.Quantity,
		)
		if err != nil {
			return nil, err
		}

		inventory = append(inventory, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inventory, nil
}

func (r *PostgresRepository) GetCoinHistory(ctx context.Context) {}
