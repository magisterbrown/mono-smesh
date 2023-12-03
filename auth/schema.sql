CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(128) UNIQUE NOT NULL,
    password VARCHAR(128) NOT NULL
);

CREATE TABLE session_token (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES players(id),
    token VARCHAR(128) NOT NULL
)
