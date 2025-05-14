CREATE TABLE balances (
    user_id INT PRIMARY KEY,
    amount DECIMAL(19, 4) NOT NULL DEFAULT 0.00 CHECK (amount >= 0),
    last_updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_balances_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_balances_last_updated_at ON balances(last_updated_at);