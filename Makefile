# Makefile

run:
	go get -d
	go run *.go

build:
	go get -d
	go build -o out.bin

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
		./...