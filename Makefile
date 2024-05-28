integration_test:
	go test ./cmd/... -count=1

unit_test:
	go test ./intenal/... -count=1

test: integration_test unit_test

test_100_times:
	.\run_test.bat

build:
	go build -o fs ./cmd/main.go

run:
	./vfs
