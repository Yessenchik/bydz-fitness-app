# BYDZ Fitness App

Infrastructure:

* PostgreSQL
* Redis
* NATS JetStream
* Prometheus
* Grafana
* Alertmanager

---


# Technologies

## Backend

* Go
* gRPC
* PostgreSQL
* Redis
* NATS JetStream
* Prometheus
* Grafana
* Alertmanager
* Docker
* Zap Logger

## Frontend (planned)

* React + TypeScript

---

# Running The Project


# Start Infrastructure

```bash
cd deploy
docker compose --env-file ../.env.docker up -d --build
```

Check containers:

```bash
docker ps
```

Expected containers:

* bydz-auth-service
* bydz-membership-service
* bydz-postgres
* bydz-redis
* bydz-nats
* bydz-prometheus
* bydz-grafana
* bydz-alertmanager

---

# Auth Service

## Features

* User registration
* Login
* JWT authentication
* Refresh tokens
* Logout
* Redis session storage
* gRPC API
* Prometheus metrics
* Structured logging
* NATS events

---

# Auth Service Endpoints

## Register

```bash
grpcurl -plaintext \
-d '{
  "email":"yessen@gmail.com",
  "password":"123456",
  "firstName":"Yessen",
  "lastName":"Zhumagali",
  "phone":"+77777777777"
}' \
localhost:50061 \
auth.v1.AuthService/Register
```

---

## Login

```bash
grpcurl -plaintext \
-d '{
  "email":"yessen@gmail.com",
  "password":"123456"
}' \
localhost:50061 \
auth.v1.AuthService/Login
```

---

## Refresh Token

```bash
grpcurl -plaintext \
-d '{
  "refreshToken":"YOUR_REFRESH_TOKEN"
}' \
localhost:50061 \
auth.v1.AuthService/RefreshToken
```

---

## Logout

```bash
grpcurl -plaintext \
-d '{
  "refreshToken":"YOUR_REFRESH_TOKEN"
}' \
localhost:50061 \
auth.v1.AuthService/Logout
```

---

## Get Profile

```bash
grpcurl -plaintext \
-H "authorization: Bearer YOUR_ACCESS_TOKEN" \
-d '{
  "userId":"USER_ID"
}' \
localhost:50061 \
auth.v1.AuthService/GetProfile
```

---

# Auth Service Events

## Published Events

### user.registered

Payload:

```json
{
  "user_id":"uuid",
  "email":"yessen@gmail.com",
  "first_name":"Yessen"
}
```

---

# Membership Service

## Features

* Membership plans
* Subscriptions
* Membership validation
* Membership cancellation
* gRPC API
* NATS events
* Prometheus metrics
* Structured logging

---

# Membership Service Ports

| Port  | Description        |
| ----- | ------------------ |
| 50062 | gRPC               |
| 9102  | Prometheus Metrics |

---

# Membership Service Endpoints

## Create Plan

```bash
grpcurl -plaintext \
-d '{
  "name":"Monthly Membership",
  "durationDays":30,
  "priceKzt":32990
}' \
localhost:50062 \
membership.v1.MembershipService/CreatePlan
```

---

## List Plans

```bash
grpcurl -plaintext \
-d '{"onlyActive":true}' \
localhost:50062 \
membership.v1.MembershipService/ListPlans
```

---

## Create Subscription

```bash
grpcurl -plaintext \
-d '{
  "userId":"USER_ID",
  "planId":"PLAN_ID"
}' \
localhost:50062 \
membership.v1.MembershipService/CreateSubscription
```

---

## Validate Membership

```bash
grpcurl -plaintext \
-d '{
  "userId":"USER_ID"
}' \
localhost:50062 \
membership.v1.MembershipService/ValidateMembership
```

---

## Cancel Subscription

```bash
grpcurl -plaintext \
-d '{
  "subscriptionId":"SUBSCRIPTION_ID"
}' \
localhost:50062 \
membership.v1.MembershipService/CancelSubscription
```

---

# Membership Service Events

## Published Events

### membership.subscription_created

```json
{
  "subscription_id":"uuid",
  "user_id":"uuid",
  "plan_id":"uuid",
  "status":"active"
}
```

---

### membership.subscription_cancelled

```json
{
  "subscription_id":"uuid",
  "user_id":"uuid",
  "plan_id":"uuid",
  "status":"cancelled"
}
```

---

# Monitoring

## Prometheus

Open:

```text
http://localhost:54303
```

---

## Grafana

Open:

```text
http://localhost:54328
```

Default credentials:

```text
admin
admin
```

---

# Metrics Endpoints

## Auth Service

```text
http://localhost:9101/metrics
```

## Membership Service

```text
http://localhost:9102/metrics
```

---

# NATS JetStream

## Open NATS Monitoring

```text
http://localhost:8222
```

---

# Test NATS Events

```bash
docker run --rm -it --network deploy_default natsio/nats-box sh
```

Inside container:

```bash
nats sub ">" -s nats://bydz-nats:4222
```

---

# Logging

Services use Zap structured logging.

Example:

```json
{
  "level":"info",
  "msg":"grpc request",
  "method":"/auth.v1.AuthService/Login",
  "status":"OK"
}
```

---

# Future Services

Planned microservices:

* notification-service
* payment-service
* attendance-service
* api-gateway
* shop-service

---

# Planned Features

* SMTP email notifications
* QR attendance
* Redis caching
* Role-based access
* API Gateway
* React frontend auth flow
* Payment integration
* Durable JetStream consumers
* Idempotent event processing

---

# Complete gRPC Testing Commands

## Auth Service Reflection

```bash
grpcurl -plaintext localhost:50061 list
```

---

## Auth Service Methods

```bash
grpcurl -plaintext localhost:50061 list auth.v1.AuthService
```

---

## Register User

```bash
grpcurl -plaintext \
-d '{
  "email":"yessen@gmail.com",
  "password":"123456",
  "firstName":"Yessen",
  "lastName":"Zhumagali",
  "phone":"+77777777777"
}' \
localhost:50061 \
auth.v1.AuthService/Register
```

---

## Login User

```bash
grpcurl -plaintext \
-d '{
  "email":"yessen@gmail.com",
  "password":"123456"
}' \
localhost:50061 \
auth.v1.AuthService/Login
```

---

## Refresh Token

```bash
grpcurl -plaintext \
-d '{
  "refreshToken":"YOUR_REFRESH_TOKEN"
}' \
localhost:50061 \
auth.v1.AuthService/RefreshToken
```

---

## Logout

```bash
grpcurl -plaintext \
-d '{
  "refreshToken":"YOUR_REFRESH_TOKEN"
}' \
localhost:50061 \
auth.v1.AuthService/Logout
```

---

## Get Profile

```bash
grpcurl -plaintext \
-H "authorization: Bearer YOUR_ACCESS_TOKEN" \
-d '{
  "userId":"USER_ID"
}' \
localhost:50061 \
auth.v1.AuthService/GetProfile
```

---

# Membership Service Testing

## Membership Service Reflection

```bash
grpcurl -plaintext localhost:50062 list
```

---

## Membership Service Methods

```bash
grpcurl -plaintext localhost:50062 list membership.v1.MembershipService
```

---

## Create Membership Plan

```bash
grpcurl -plaintext \
-d '{
  "name":"Monthly Membership",
  "durationDays":30,
  "priceKzt":32990
}' \
localhost:50062 \
membership.v1.MembershipService/CreatePlan
```

---

## List Membership Plans

```bash
grpcurl -plaintext \
-d '{
  "onlyActive":true
}' \
localhost:50062 \
membership.v1.MembershipService/ListPlans
```

---

## Create Subscription

```bash
grpcurl -plaintext \
-d '{
  "userId":"USER_ID",
  "planId":"PLAN_ID"
}' \
localhost:50062 \
membership.v1.MembershipService/CreateSubscription
```

---

## Get Subscription

```bash
grpcurl -plaintext \
-d '{
  "subscriptionId":"SUBSCRIPTION_ID"
}' \
localhost:50062 \
membership.v1.MembershipService/GetSubscription
```

---

## Validate Membership

```bash
grpcurl -plaintext \
-d '{
  "userId":"USER_ID"
}' \
localhost:50062 \
membership.v1.MembershipService/ValidateMembership
```

---

## Cancel Subscription

```bash
grpcurl -plaintext \
-d '{
  "subscriptionId":"SUBSCRIPTION_ID"
}' \
localhost:50062 \
membership.v1.MembershipService/CancelSubscription
```

---

## List User Subscriptions

```bash
grpcurl -plaintext \
-d '{
  "userId":"USER_ID"
}' \
localhost:50062 \
membership.v1.MembershipService/ListUserSubscriptions
```

---

# Metrics Testing

## Auth Service Metrics

```bash
curl localhost:9101/metrics
```

---

## Membership Service Metrics

```bash
curl localhost:9102/metrics
```

---

# NATS Testing

## Open NATS Box

```bash
docker run --rm -it --network deploy_default natsio/nats-box sh
```

---

## Subscribe To All Events

```bash
nats sub ">" -s nats://bydz-nats:4222
```

---

## Subscribe To Auth Events

```bash
nats sub "user.registered" -s nats://bydz-nats:4222
```

---

## Subscribe To Membership Events

```bash
nats sub "membership.>" -s nats://bydz-nats:4222
```

---

# Docker Commands

## Start Everything

```bash
docker compose --env-file .env.docker -f deploy/docker-compose.yml up -d
```

---

## Restart Auth Service

```bash
docker restart bydz-auth-service
```

---

## Restart Membership Service

```bash
docker restart bydz-membership-service
```

---

## Service Logs

```bash
docker logs -f bydz-auth-service
```

```bash
docker logs -f bydz-membership-service
```

---

## Running Containers

```bash
docker ps
```

---

# Team

* Yessen
* Zhan
* Bekbauly
* Danial
