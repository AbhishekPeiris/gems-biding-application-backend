CREATE TABLE IF NOT EXISTS auctions (
    id BIGSERIAL PRIMARY KEY,
    gem_id BIGINT NOT NULL REFERENCES gems(id) ON DELETE CASCADE,
    start_price NUMERIC(15,2) NOT NULL,
    current_price NUMERIC(15,2) NOT NULL,
    min_increment NUMERIC(15,2) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('SCHEDULED','LIVE','ENDED')),
    winner_id BIGINT REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
