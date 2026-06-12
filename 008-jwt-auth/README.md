# Go JWT Auth API

JWT-based authentication. Role-based access (ADMIN/USER). Net/http + PostgreSQL.

## Start

```bash
docker-compose up -d
go run main.go
```

Server: `http://localhost:8080` | DB: `localhost:5432`

---

## Curl Commands

### Signup User
```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "1234567890",
    "password": "password123",
    "role": "USER"
  }'
```

### Signup Admin
```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Admin",
    "last_name": "User",
    "email": "admin@example.com",
    "phone": "0987654321",
    "password": "password123",
    "role": "ADMIN"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```
Save the `token` from response.

### Get All Users (ADMIN only)
```bash
TOKEN="your-token-here"
curl -X GET http://localhost:8080/users \
  -H "Authorization: Bearer $TOKEN"
```

### Get User by ID (ADMIN or self)
```bash
TOKEN="your-token-here"
USER_ID="user-uuid-here"
curl -X GET http://localhost:8080/user/$USER_ID \
  -H "Authorization: Bearer $TOKEN"
```

---

## Routes Summary

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| POST | `/signup` | None | Create user (role: ADMIN or USER) |
| POST | `/login` | None | Login, returns token + refresh_token |
| GET | `/users` | Bearer + ADMIN | List all users |
| GET | `/user/{id}` | Bearer | Get user (ADMIN or self) |

## Database

Auto-creates schema on first `docker-compose up`.

Default: `postgres:postgres` @ `localhost:5432` (db: `usersdb`)

## Cleanup

```bash
docker-compose down      # Stop
docker-compose down -v   # Stop + wipe data
```
