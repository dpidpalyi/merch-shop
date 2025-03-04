package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"merch-shop/internal/models"
)

type Repository interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Add(ctx context.Context, u *models.User) error
	BuyItem(ctx context.Context, userID int, itemName string) error
	SendCoin(ctx context.Context, senderID, receiverID int, amount int) error
	GetBalance(ctx context.Context, userID int) (int, error)
	GetInventory(ctx context.Context, userID int) ([]*models.InventoryItem, error)
	GetCoinHistory(ctx context.Context, userID int) (*models.CoinHistory, error)
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
	    FROM active_users
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
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	    INSERT INTO users(username, password_hash)
	    VALUES ($1, $2)
	    RETURNING id, created_at`

	args := []any{u.Username, u.PasswordHash}

	err = tx.QueryRowContext(ctx, query, args...).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `
	    INSERT INTO coins(user_id)
	    VALUES ($1)`

	_, err = tx.ExecContext(ctx, query, u.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetItemByName(ctx context.Context, itemName string) (*models.Item, error) {
	query := `
	    SELECT id, price
	    FROM items
	    WHERE type = $1`

	item := &models.Item{
		Name: itemName,
	}

	err := r.DB.QueryRowContext(ctx, query, itemName).Scan(
		&item.ID,
		&item.Price,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return item, nil
}

func (r *PostgresRepository) BuyItem(ctx context.Context, userID int, itemName string) error {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	     SELECT balance
	     FROM coins
	     JOIN active_users ON coins.user_id = active_users.id
	     WHERE user_id = $1 FOR UPDATE`

	var balance int
	err = tx.QueryRowContext(ctx, query, userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user: %w", ErrRecordNotFound)
		}
		return err
	}

	query = `
	    SELECT id, price
	    FROM item
	    WHERE type = $1`

	var item models.Item

	err = tx.QueryRowContext(ctx, query, itemName).Scan(
		&item.ID,
		&item.Price,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("item: %w", ErrRecordNotFound)
		}
		return err
	}

	if err := checkBalance(balance, item.Price); err != nil {
		return err
	}

	query = `
	    INSERT INTO inventory(user_id, item_id)
	    VALUES ($1, $2)
	    ON CONFLICT (user_id, item_id)
	    DO UPDATE SET quantity = inventory.quantity + 1`

	args := []any{userID, item.ID}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	query = `
	    UPDATE coins
	    SET balance = balance - $2
	    WHERE user_id = $1`

	args = []any{userID, item.Price}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresRepository) SendCoin(ctx context.Context, senderID, receiverID int, amount int) error {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	     SELECT balance
	     FROM coins
	     JOIN active_users ON coins.user_id = active_users.id
	     WHERE user_id = $1 FOR UPDATE`

	var balance int
	err = tx.QueryRowContext(ctx, query, senderID).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("sender user: %w", ErrRecordNotFound)
		}
		return err
	}

	if err := checkBalance(balance, amount); err != nil {
		return err
	}

	query = `
	     INSERT INTO transaction(sender_id, receiver_id, amount)
	     VALUES ($1, $2, $3)`

	args := []any{senderID, receiverID, amount}

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

	return tx.Commit()
}

func (r *PostgresRepository) GetBalance(ctx context.Context, userID int) (int, error) {
	query := `
	    SELECT balance
	    FROM coins
	    JOIN active_users ON coins.user_id = active_users.id
	    WHERE user_id = $1`

	var balance int

	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user: %w", ErrRecordNotFound)
		}
		return 0, err
	}

	return balance, nil
}

func (r *PostgresRepository) GetInventory(ctx context.Context, userID int) ([]*models.InventoryItem, error) {
	query := `
	    SELECT it.type, i.quantity
	    FROM inventory AS i
	    JOIN item AS it ON i.item_id = it.id
	    WHERE user_id = $1`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventory := []*models.InventoryItem{}

	for rows.Next() {
		var item models.InventoryItem
		err := rows.Scan(
			&item.Type,
			&item.Quantity,
		)
		if err != nil {
			return nil, err
		}

		inventory = append(inventory, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inventory, nil
}

func (r *PostgresRepository) GetCoinHistory(ctx context.Context, userID int) (*models.CoinHistory, error) {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	    SELECT u1.username, t.amount
	    FROM transaction AS t
	    JOIN users AS u1 ON t.sender_id = u1.id
	    JOIN users AS u2 ON t.receiver_id = u2.id
	    WHERE t.receiver_id = $1`

	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	received := []*models.CoinTransaction{}

	for rows.Next() {
		var r models.CoinTransaction
		err := rows.Scan(
			&r.FromUser,
			&r.Amount,
		)
		if err != nil {
			return nil, err
		}

		received = append(received, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	query = `
	    SELECT u2.username, t.amount
	    FROM transaction AS t
	    JOIN users AS u1 ON t.sender_id = u1.id
	    JOIN users AS u2 ON t.receiver_id = u2.id
	    WHERE t.sender_id = $1`

	rows, err = tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sent := []*models.CoinTransaction{}

	for rows.Next() {
		var s models.CoinTransaction
		err := rows.Scan(
			&s.ToUser,
			&s.Amount,
		)
		if err != nil {
			return nil, err
		}

		sent = append(sent, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	coinHistory := &models.CoinHistory{
		Received: received,
		Sent:     sent,
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return coinHistory, nil
}
