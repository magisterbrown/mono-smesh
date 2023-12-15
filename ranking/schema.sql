DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(128) UNIQUE NOT NULL
);

CREATE TABLE submissions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    container_id VARCHAR(128) NOT NULL,
    raiting FLOAT NOT NULL,
    sigma FLOAT NOT NULL,
    broken INT DEFAULT 0,
    FOREIGN KEY(user_id) REFERENCES players(id)
);

INSERT INTO players (user_name)
VALUES ('JohnDoe'), ('goof'), ('sam'), ('biilyherington');
