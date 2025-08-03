-- Users table
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- File uploads table
CREATE TABLE file_uploads (
                              id SERIAL PRIMARY KEY,
                              filename VARCHAR(255) NOT NULL,
                              content_type VARCHAR(100) NOT NULL,
                              size BIGINT NOT NULL,
                              file_path VARCHAR(500) NOT NULL,
                              user_agent TEXT,
                              remote_addr VARCHAR(45),
                              user_id INTEGER REFERENCES users(id),
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Revoked tokens table (optional, for persistent token revocation)
CREATE TABLE revoked_tokens (
                                id SERIAL PRIMARY KEY,
                                token TEXT NOT NULL,
                                revoked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);