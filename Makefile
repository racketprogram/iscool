integration_test:
	go test ./cmd/... -coverprofile=integration_coverage  -count=1
	go tool cover -func=integration_coverage

unit_test:
	go test  ./internal/... -coverprofile=unit_coverage -count=1
	go tool cover -func=unit_coverage

test: integration_test unit_test

test_100_times:
	.\run_test.bat

build:
	go build -o fs ./cmd/main.go

run:
	./vfs
