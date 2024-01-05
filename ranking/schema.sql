/* DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
*/

CREATE TABLE IF NOT EXISTS players (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(128) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS submissions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    file_name VARCHAR(128) NOT NULL,
    container_id VARCHAR(128) NOT NULL,
    raiting FLOAT NOT NULL,
    sigma FLOAT NOT NULL,
    broken BOOL DEFAULT 'f',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES players(id)
);

CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    recording json
);

CREATE TABLE IF NOT EXISTS seating (
    id SERIAL PRIMARY KEY,
    match_id INT NOT NULL,
    submission_id INT NOT NULL,
    change FLOAT NOT NULL,
    spot VARCHAR(128) NOT NULL,
    FOREIGN KEY(match_id) REFERENCES matches(id),
    FOREIGN KEY(submission_id) REFERENCES submissions(id)
);

INSERT INTO players (user_name)
VALUES ('JohnDoe'), ('goof'), ('sam'), ('biilyherington') ON CONFLICT DO NOTHING;
