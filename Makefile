build:
	@go build -o bin/hotelreservation

run: build
	@./bin/hotelreservation

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...