CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   username VARCHAR(50) UNIQUE NOT NULL,
   password VARCHAR(50) NOT NULL,
);

CREATE TABLE notes (
   id SERIAL PRIMARY KEY,
   user_id INT REFERENCES users(id),
   title VARCHAR(100) NOT NULL,
   content TEXT NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


INSERT INTO users (username, password) VALUES
    ('testuser1', '1234');

INSERT INTO users (username, password) VALUES
    ('testuser2', '12345');

INSERT INTO users (username, password) VALUES
    ('testuser3', '123456');

INSERT INTO users (username, password) VALUES
    ('testuser4', '1234567');

INSERT INTO users (username, password) VALUES
    ('testuser5', '12345678');

INSERT INTO users (username, password) VALUES
    ('testuser6', '123456789');