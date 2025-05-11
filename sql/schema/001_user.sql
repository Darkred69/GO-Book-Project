-- goose postgres://postgres:{username}@{database_IP}:{database_port}/{databasename}?sslmode=disable
--+goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),  
    name TEXT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE users;