CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_user_id INT,
    to_user_id INT,
    amount DECIMAL(19, 4) NOT NULL CHECK (amount > 0),
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_from_user
        FOREIGN KEY(from_user_id)
        REFERENCES users(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_to_user
        FOREIGN KEY(to_user_id)
        REFERENCES users(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_transactions_from_user_id ON transactions(from_user_id);
CREATE INDEX idx_transactions_to_user_id ON transactions(to_user_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);