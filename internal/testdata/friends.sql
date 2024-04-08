-- testdata/friends.sql
DELETE FROM friend_connections;
DELETE FROM users;

-- Sample users
INSERT INTO users (id, name, email) VALUES (1, 'John Doe', 'john@example.com');
INSERT INTO users (id, name, email) VALUES (2, 'Jane Doe', 'jane@example.com');
INSERT INTO users (id, name, email) VALUES (3, 'Alice Smith', 'alice@example.com');
INSERT INTO users (id, name, email) VALUES (4, 'Bob Johnson', 'bob@example.com');

-- Sample friend connections
INSERT INTO friend_connections (user_id1, user_id2) VALUES (1, 2); -- John and Jane are friends
