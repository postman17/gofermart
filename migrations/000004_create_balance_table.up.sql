CREATE TABLE IF NOT EXISTS user_balances (
    user_id INT PRIMARY KEY,
    current NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    withdrawn NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_balance_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE,
        
    CONSTRAINT chk_current_positive CHECK (current >= 0.00),
    CONSTRAINT chk_withdrawn_positive CHECK (withdrawn >= 0.00)
);
