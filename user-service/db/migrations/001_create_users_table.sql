CREATE TABLE users (
  user_id SERIAL PRIMARY KEY,
  name TEXT NOT NULL
);

INSERT INTO users (name) VALUES
('Hasan'),
('Ali'),
('Fatima');
