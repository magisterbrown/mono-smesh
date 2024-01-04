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
    broken INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES players(id)
);

CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    subm_1 INT NOT NULL,
    subm_2 INT NOT NULL,
    subm_1_change FLOAT NOT NULL,
    subm_2_change FLOAT NOT NULL,
    recording json,
    FOREIGN KEY(subm_1) REFERENCES submissions(id),
    FOREIGN KEY(subm_2) REFERENCES submissions(id)

);

INSERT INTO players (user_name)
VALUES ('JohnDoe'), ('goof'), ('sam'), ('biilyherington') ON CONFLICT DO NOTHING;
