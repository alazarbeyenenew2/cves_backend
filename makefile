migrate-down:
	- migrate -database postgresql://test_user:secret_password@localhost:5432/cves?sslmode=disable -path internal/constant/db/schemas -verbose down $(N)
migrate-up:
	- migrate -database postgresql://test_user:secret_password@localhost:5432/cves?sslmode=disable -path internal/constant/db/schemas -verbose up
migrate-create:
	- migrate create -ext sql -dir internal/constant/db/schemas -tz "UTC" $(name)
swagger:
	-swag fmt && swag init -g cmd/main.go
run:
	go run cmd/main.go
sqlc:
	cd ./config && sqlc generate
