-- a simple table of super heroes
CREATE TABLE IF NOT EXISTS superhero (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    power smallint NOT NULl,
    melee boolean NOT NULL
)