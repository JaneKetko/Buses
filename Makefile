# Makefile

run: 
	./cmd/service/service & ./cmd/client/client

run-service:
	./cmd/service/service

run-client:
	./cmd/client/client -c :8002

build_service:
	CGO_ENABLED=0 go build -o ./cmd/service/service github.com/JaneKetko/Buses/cmd/service

build_client:
	CGO_ENABLED=0 go build -o ./cmd/client/client github.com/JaneKetko/Buses/cmd/client

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

