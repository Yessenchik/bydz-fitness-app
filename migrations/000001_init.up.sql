CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS membership_plans (
                                                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    name VARCHAR(100) NOT NULL,
    duration_months INT NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    description TEXT,

    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_membership_duration
    CHECK (duration_months IN (1, 6, 12)),

    CONSTRAINT chk_membership_price
    CHECK (price >= 0)
    );

CREATE TABLE IF NOT EXISTS subscriptions (
                                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL,
    membership_plan_id UUID NOT NULL,

    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,

                           is_active BOOLEAN NOT NULL DEFAULT TRUE,
                           status VARCHAR(20) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_subscriptions_membership_plan
    FOREIGN KEY (membership_plan_id)
    REFERENCES membership_plans(id)
                       ON DELETE RESTRICT,

    CONSTRAINT chk_subscription_status
    CHECK (status IN ('ACTIVE', 'EXPIRED', 'CANCELLED')),

    CONSTRAINT chk_subscription_dates
    CHECK (end_date > start_date)
    );

CREATE TABLE IF NOT EXISTS attendances (
                                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL,

    check_in_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    date DATE NOT NULL DEFAULT CURRENT_DATE,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_membership_plans_is_deleted
    ON membership_plans(is_deleted);

CREATE INDEX IF NOT EXISTS idx_membership_plans_duration_months
    ON membership_plans(duration_months);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id
    ON subscriptions(user_id);

CREATE INDEX IF NOT EXISTS idx_subscriptions_membership_plan_id
    ON subscriptions(membership_plan_id);

CREATE INDEX IF NOT EXISTS idx_subscriptions_status
    ON subscriptions(status);

CREATE INDEX IF NOT EXISTS idx_subscriptions_is_active
    ON subscriptions(is_active);

CREATE INDEX IF NOT EXISTS idx_subscriptions_end_date
    ON subscriptions(end_date);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_active
    ON subscriptions(user_id, is_active, status);

CREATE INDEX IF NOT EXISTS idx_attendances_user_id
    ON attendances(user_id);

CREATE INDEX IF NOT EXISTS idx_attendances_date
    ON attendances(date);

CREATE INDEX IF NOT EXISTS idx_attendances_check_in_time
    ON attendances(check_in_time);

CREATE INDEX IF NOT EXISTS idx_attendances_user_date
    ON attendances(user_id, date);

CREATE INDEX IF NOT EXISTS idx_attendances_user_check_in_time
    ON attendances(user_id, check_in_time);