CREATE TABLE song_details (
    seq_num SERIAL PRIMARY KEY,
    id INT,
    title VARCHAR(255),
    release_date DATE,
    artist VARCHAR(255),
    lyrics TEXT,
    link VARCHAR(255)
);
