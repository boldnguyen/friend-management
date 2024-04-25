-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create friend_connections table
CREATE TABLE friend_connections (
    id SERIAL PRIMARY KEY,
    user_id1 INT REFERENCES users(id) NOT NULL,
    user_id2 INT REFERENCES users(id) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (user_id1, user_id2)
);

-- Create subscriptions table
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    requestor VARCHAR(255) NOT NULL,
    target VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()  
);
