CREATE TABLE IF NOT EXISTS tracked_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_email TEXT NOT NULL,
    merchant TEXT NOT NULL,
    category TEXT,
    amount NUMERIC(10,2),
    raw_text TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('active', 'canceled')),
    provider TEXT NOT NULL DEFAULT 'direct',
    provider_ref TEXT,
    email_token TEXT,
    confidence NUMERIC(3,2),
    match_type TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    canceled_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_tracked_email ON tracked_subscriptions(user_email);
CREATE INDEX IF NOT EXISTS idx_tracked_status ON tracked_subscriptions(status);