CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    name CHAR(100) NOT NULL,
    age INT NOT NULL
);
CREATE TABLE IF NOT EXISTS friends(
                                      id SERIAL PRIMARY KEY,
                                      user1_id INT,
                                      user2_id INT
);