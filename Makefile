# Makefile

run: 
	./bin/service & ./bin/client

run-service:
	go run cmd/service/main.go

run-client:
	go run cmd/client/main.go

build_service:
	go build -o ./bin/service github.com/JaneKetko/Buses/cmd/service

build_client:
	go build -o ./bin/client github.com/JaneKetko/Buses/cmd/client

build: build_service build_client

test:
	go test -v ./... -race -coverprofile=coverage.txt -covermode=atomic GOCACHE=off

lint:
	golangci-lint run \
		--exclude-use-default=false \
		--enable=golint \
		--enable=gocyclo \
		--enable=goconst \
		--enable=unconvert \
		--enable=dupl \
		--enable=maligned \
		--enable=depguard \
		--enable=misspell \
		--enable=unparam \
		--enable=nakedret \
		--enable=prealloc \
		--enable=scopelint \
		--enable=gochecknoglobals \
		--enable=gochecknoinits \
		--gocyclo.min-complexity 10 \
		./...

