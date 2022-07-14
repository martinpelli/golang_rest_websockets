DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id VARCHAR(32) PRIMARY KEY,
    password VARCHAR(25) NOT NULL,
    email VARCHAR(25) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
)