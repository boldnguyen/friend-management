-- friends.sql
INSERT INTO users (id, name, email) VALUES
(1, 'Andy', 'andy@example.com'), -- Thêm người dùng user1
(2, 'John', 'john@example.com') -- Add Bob as a new user

ON CONFLICT DO NOTHING; -- Skip inserting if the user already exists

INSERT INTO friend_connections (user_id1, user_id2) VALUES
(1, 2)
ON CONFLICT DO NOTHING; -- Skip inserting if the friend connection already exists
