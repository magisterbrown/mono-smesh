CREATE TABLE players (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name VARCHAR(128) UNIQUE NOT NULL,
    password VARCHAR(128) NOT NULL
);

CREATE TABLE submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INT NOT NULL,
    container_id VARCHAR(128) NOT NULL,
    raiting FLOAT NOT NULL,
    sigma FLOAT NOT NULL,
    broken INTEGER DEFAULT 0,
    FOREIGN KEY(user_id) REFERENCES players(id)
);

CREATE TABLE sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INT NOT NULL,
    token VARCHAR(128) UNIQUE NOT NULL,
    FOREIGN KEY(user_id) REFERENCES players(id)
);
INSERT INTO players (user_name, password) VALUES 
    ("aa", "140.0"),
    ("bb", "130.0"),
    ("bc", "130.0");

INSERT INTO sessions (user_id, token) VALUES
    (1, "tiktok");

