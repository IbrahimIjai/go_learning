# sqlc tutorial

## Commands

```bash
# 1. Start Postgres (Docker)
docker start sqlc-tutorial-pg

# 2. Apply migrations
~/go/bin/migrate -database "postgres://tutorial:tutorial@localhost:5432/tutorial?sslmode=disable" -path migrations up

# 3. Regenerate Go code after editing query.sql / migrations
sqlc generate

# 4. Run the app
go run .
```

Roll back the last migration with `... -path migrations down 1`.
