# Makefile

run: 
	./bin/server & ./bin/client

run-service:
	go run cmd/service/main.go

run-grpcserver:
	go run cmd/grpcserver/main.go

run-client:
	go run cmd/client/main.go

build_server:
	go build -o ./bin/server github.com/JaneKetko/Buses/cmd/grpcserver

build_client:
	go build -o ./bin/client github.com/JaneKetko/Buses/cmd/client

build: build_server build_client

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

