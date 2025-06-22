CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,
    username TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS trackings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    tracking_number TEXT NOT NULL,
    carrier_code TEXT,
    status TEXT,
    last_update TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
); 