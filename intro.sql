CREATE TABLE IF NOT EXISTS Albums (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    price FLOAT NOT NULL
);

INSERT INTO Albums (title, artist, price)
VALUES ('imagine', 'john lenon',19);