-- subscriptions.sql
INSERT INTO subscriptions (requestor, target) VALUES 
('andy@example.com','john@example.com') -- user1 follows user2
ON CONFLICT (requestor, target) DO NOTHING;