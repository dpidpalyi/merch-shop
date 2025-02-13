CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	username VARCHAR(50) UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS coins (
	user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
	balance INT NOT NULL DEFAULT 1000
);

CREATE INDEX idx_user_id ON coins(user_id);

CREATE TABLE IF NOT EXISTS item (
	id SERIAL PRIMARY KEY,
	type VARCHAR(50) UNIQUE NOT NULL,
	price INT NOT NULL
);

CREATE TABLE IF NOT EXISTS inventory (
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id) ON DELETE CASCADE,
	item_id INT REFERENCES item(id) ON DELETE CASCADE,
	quantity INT NOT NULL DEFAULT 1,
	UNIQUE (user_id, item_id)
);

CREATE TABLE IF NOT EXISTS transaction (
	id SERIAL PRIMARY KEY,
	sender_id INT REFERENCES users(id) ON DELETE SET NULL,
	receiver_id INT REFERENCES users(id) ON DELETE SET NULL,
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

--CREATE TABLE IF NOT EXISTS auth_token (
--	id SERIAL PRIMARY KEY,
--	user_id INT REFERENCES users(id) ON DELETE CASCADE,
--	token TEXT NOT NULL UNIQUE,
--	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
--);
