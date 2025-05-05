CREATE TABLE rides (
  ride_id SERIAL PRIMARY KEY,
  source TEXT NOT NULL,
  destination TEXT NOT NULL,
  distance INT NOT NULL,
  cost INT NOT NULL
);

INSERT INTO rides (source, destination, distance, cost) VALUES
('Karachi', 'Lahore', 1200, 5000),
('Islamabad', 'Peshawar', 200, 1500),
('Multan', 'Faisalabad', 400, 2500);
