CREATE TABLE items (
                       id SERIAL PRIMARY KEY,
                       parent_id INT REFERENCES items(id) ON DELETE CASCADE,
                       name VARCHAR(255) NOT NULL,
                       is_directory BOOLEAN NOT NULL,
                       file_path TEXT,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);