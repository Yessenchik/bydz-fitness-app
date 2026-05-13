## How to Run the Project
1. Build the project locally
```text
go build ./cmd/server
```
2. Run with Docker Compose
```text
docker compose up --build app
```
3. Run all infrastructure
```text
docker compose up --build
```
This starts:
```text
membership-service
postgres
redis
nats
prometheus
grafana
```
## PostgreSQL Check

Enter PostgreSQL container:
```text
docker exec -it membership-postgres psql -U gym -d membership_db
```
Show tables and exit:
```text
\dt
\q
```
## Testing gRPC in Postman

Open Postman and create a new gRPC request.

Address:
```text
localhost:50051
```
Example: CreateMembershipPlan
```text
{
"name": "1 Month Plan",
"duration_months": 1,
"price": 15000,
"description": "Basic monthly gym access"
}
```
Example: ListMembershipPlans
```text
{
  "limit": 10,
  "offset": 0
}
```
Example: CreateSubscription
```text
{
  "user_id": "11111111-1111-1111-1111-111111111111",
  "membership_plan_id": "PLAN_ID_HERE",
  "user_email": "test@example.com"
}
```
Example: CheckAccess, RegisterAttendance, GetMonthlyAttendanceStats
```text
{
  "user_id": "11111111-1111-1111-1111-111111111111"
}
```
## JWT Middleware

The service contains JWT middleware.
```text
user_id
role
email
```
from JWT claims and adds them to context.

Admin-only operations:
```text
CreateMembershipPlan
UpdateMembershipPlan
DeleteMembershipPlan
GetGlobalStats
```

## Running Tests
Unit Tests

Unit tests do not require Docker.
```text
go test ./tests/unit -v
```
Unit tests cover:
```text
subscription creation
access checking
subscription cancellation
attendance registration
monthly statistics
```
Integration Tests

Integration tests require PostgreSQL.
```text
$env:TEST_DATABASE_DSN="postgres://gym:gym_password@127.0.0.1:5433/membership_db?sslmode=disable"; go test ./tests/integration -v
```
Integration tests cover:
```text
attendance repository
subscription repository
gRPC handler
real PostgreSQL connection
```
Run All Tests
```text
$env:TEST_DATABASE_DSN="postgres://gym:gym_password@127.0.0.1:5433/membership_db?sslmode=disable"; go test ./... -v
```
```text
docker compose ps
docker compose logs app
docker compose down
docker compose down -v
```

## Current Project Status
```text
Go build                      PASS
Docker build                  PASS
PostgreSQL connection          PASS
Redis connection               PASS
NATS connection                PASS
gRPC server startup            PASS
PostgreSQL migrations          PASS
Postman gRPC testing           PASS
Unit tests                     PASS
Integration tests              PASS
```