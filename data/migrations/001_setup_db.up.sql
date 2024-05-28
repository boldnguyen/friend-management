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
    user_id1 INT NOT NULL,
    user_id2 INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (user_id1, user_id2),
    FOREIGN KEY (user_id1) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id2) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (user_id1 <> user_id2) -- Ensures a user cannot friend themselves
);

-- Create subscriptions table
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    requestor VARCHAR(255) NOT NULL REFERENCES users(email),
    target VARCHAR(255) NOT NULL REFERENCES users(email),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (requestor, target)
);

-- Create blocks table
CREATE TABLE blocks (
    id SERIAL PRIMARY KEY,
    requestor INT NOT NULL,
    target INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (requestor, target),
    FOREIGN KEY (requestor) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (target) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (requestor <> target) -- Ensures a user cannot block themselves
);

