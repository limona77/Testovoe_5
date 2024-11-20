CREATE TABLE songs (
                       id SERIAL PRIMARY KEY,
                       group_name TEXT NOT NULL,
                       song_name TEXT NOT NULL,
                       release_date DATE,
                       text TEXT,
                       link TEXT,
                       created_at TIMESTAMP DEFAULT NOW(),
                       updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_group_name ON songs (group_name);
CREATE INDEX idx_song_name ON songs (song_name);
CREATE INDEX idx_release_date ON songs (release_date);