-- friends.sql
INSERT INTO users (id, name, email) VALUES
(1, 'John Doe', 'john@example.com'),
(2, 'Jane Smith', 'jane@example.com'),
(3, 'Alice Johnson', 'alice@example.com'),
(4, 'Bob Brown', 'bob@example.com') -- Add Bob as a new user

ON CONFLICT DO NOTHING; -- Skip inserting if the user already exists

INSERT INTO friend_connections (user_id1, user_id2) VALUES
(1, 2),
(4, 2) 
ON CONFLICT DO NOTHING; -- Skip inserting if the friend connection already exists
