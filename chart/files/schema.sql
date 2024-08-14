CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL
);

INSERT INTO users (username, password) VALUES ('john_doe', 'password123');
INSERT INTO users (username, password) VALUES ('jane_smith', 'securePassword456');
INSERT INTO users (username, password) VALUES ('alice_johnson', 'aliceSecret789');

INSERT INTO messages (user_id, message) VALUES (1, 'message for john');
INSERT INTO messages (user_id, message) VALUES (2, 'message for jane');
INSERT INTO messages (user_id, message) VALUES (3, 'message for alice');
