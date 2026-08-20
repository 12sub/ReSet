CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_email TEXT NOT NULL,
    plan_code TEXT NOT NULL,
    paystack_sub_id TEXT NOT NULL UNIQUE,
    paystack_email_token TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('active', 'canceled', 'pending')),
    amount INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    canceled_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_status ON subscriptions(status);