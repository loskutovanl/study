CREATE TABLE IF NOT EXISTS friends(
  id SERIAL PRIMARY KEY,
  user1_id INT,
  user2_id INT
);