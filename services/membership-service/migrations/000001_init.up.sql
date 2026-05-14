CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE membership_plans (
                                  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                  name TEXT NOT NULL,
                                  duration_days INT NOT NULL,
                                  price_kzt BIGINT NOT NULL,

                                  is_active BOOLEAN NOT NULL DEFAULT true,

                                  created_at TIMESTAMP NOT NULL DEFAULT now(),
                                  updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE subscriptions (
                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                               user_id UUID NOT NULL,
                               plan_id UUID NOT NULL REFERENCES membership_plans(id),

                               status TEXT NOT NULL,

                               starts_at TIMESTAMP NOT NULL,
                               expires_at TIMESTAMP NOT NULL,

                               cancelled_at TIMESTAMP,

                               created_at TIMESTAMP NOT NULL DEFAULT now(),
                               updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_subscriptions_user_id
    ON subscriptions(user_id);

CREATE INDEX idx_subscriptions_status
    ON subscriptions(status);