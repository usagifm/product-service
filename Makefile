fmt:
	go fmt ./...

build:
	go build .
	
run:
	go run main.go 

test:
	go test ./...

migrate up:

	migrate -database "postgres://root:123456@localhost:5432/product_service?sslmode=disable" -path=db/migrations up


migrate down: 
	
	migrate -database "postgres://root:123456@localhost:5432/product_service?sslmode=disable" -path=db/migrations down

