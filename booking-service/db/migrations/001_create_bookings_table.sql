CREATE TABLE bookings (
  booking_id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  ride_id INT NOT NULL,
  time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO bookings (user_id, ride_id) VALUES
(1, 1),
(2, 2),
(3, 3);
