-- 1. QR Transactions Table (Tracks customer scans and payments)
CREATE TABLE qr_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    
    customer_display_id VARCHAR(50) NOT NULL, -- e.g., #QR-8820
    bill_amount NUMERIC(15, 2) NOT NULL,
    reward_issued NUMERIC(15, 2) NOT NULL,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Performance Indices
CREATE INDEX idx_qr_transactions_store_id_created ON qr_transactions(store_id, created_at DESC);
