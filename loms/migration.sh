goose -dir ./migrations postgres "postgres://user:password@postgres_loms:5432/loms?sslmode=disable" status
goose -dir ./migrations postgres "postgres://user:password@postgres_loms:5432/loms?sslmode=disable" up
