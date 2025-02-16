CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	username VARCHAR(50) UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_users_is_active ON users(is_active);

CREATE VIEW active_users AS
SELECT id, username, password_hash, created_at
FROM users
WHERE is_active = TRUE;

CREATE TABLE IF NOT EXISTS coins (
	user_id INT PRIMARY KEY REFERENCES users(id),
	-- Set default balance to 100 for easier tests
	balance INT NOT NULL DEFAULT 100
);

CREATE INDEX idx_user_id ON coins(user_id);

CREATE TABLE IF NOT EXISTS item (
	id SERIAL PRIMARY KEY,
	type VARCHAR(50) UNIQUE NOT NULL,
	price INT NOT NULL
);

CREATE TABLE IF NOT EXISTS inventory (
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id),
	item_id INT REFERENCES item(id) ON DELETE CASCADE,
	quantity INT NOT NULL DEFAULT 1,
	UNIQUE (user_id, item_id)
);

CREATE TABLE IF NOT EXISTS transaction (
	id SERIAL PRIMARY KEY,
	sender_id INT REFERENCES users(id),
	receiver_id INT REFERENCES users(id),
	amount INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transaction_id ON transaction(sender_id, receiver_id);

INSERT INTO item(type, price)
VALUES
       ('t-shirt', 80),
       ('cup', 20),
       ('book', 50),
       ('pen', 10),
       ('powerbank', 200),
       ('hoody', 300),
       ('umbrella', 200),
       ('socks', 10),
       ('wallet', 50),
       ('pink-hoody', 500);
