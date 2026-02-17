CREATE TABLE IF NOT EXISTS gems (
    id BIGSERIAL PRIMARY KEY,
    seller_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    carat NUMERIC(10,2) NOT NULL,
    color VARCHAR(100),
    clarity VARCHAR(100),
    origin VARCHAR(100),
    certificate VARCHAR(255),
    image_url TEXT,
    status VARCHAR(20) NOT NULL CHECK (status IN ('AVAILABLE','AUCTION','SOLD')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
