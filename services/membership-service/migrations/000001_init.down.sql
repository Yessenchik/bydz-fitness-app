DROP INDEX IF EXISTS idx_attendances_user_check_in_time;
DROP INDEX IF EXISTS idx_attendances_user_date;
DROP INDEX IF EXISTS idx_attendances_check_in_time;
DROP INDEX IF EXISTS idx_attendances_date;
DROP INDEX IF EXISTS idx_attendances_user_id;

DROP INDEX IF EXISTS idx_subscriptions_user_active;
DROP INDEX IF EXISTS idx_subscriptions_end_date;
DROP INDEX IF EXISTS idx_subscriptions_is_active;
DROP INDEX IF EXISTS idx_subscriptions_status;
DROP INDEX IF EXISTS idx_subscriptions_membership_plan_id;
DROP INDEX IF EXISTS idx_subscriptions_user_id;

DROP INDEX IF EXISTS idx_membership_plans_duration_months;
DROP INDEX IF EXISTS idx_membership_plans_is_deleted;

DROP TABLE IF EXISTS attendances;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS membership_plans;

DROP EXTENSION IF EXISTS "uuid-ossp";