# Define variables
GO_CMD=go
MAIN=./cmd/bin/main.go

test:
	go test -v ./... -cover

run:
	$(GO_CMD) run $(MAIN) serve-http

hot:
	@echo " >> Installing gin if not installed"
	@go install github.com/codegangsta/gin@latest
	@gin -i -p 9002 -a 9090 --path cmd/bin --build cmd/bin serve-http

goose-create:
# example : make goose-create name=create_users_table
	@echo " >> Installing goose if not installed"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
ifndef name
	$(error Usage: make goose-create name=<table_name>)
else
	@goose -dir db/migrations create $(name) sql
endif

goose-up:
# example : make goose-up
	@echo " >> Installing goose if not installed"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@goose -dir db/migrations postgres "host=localhost user=postgres password=21012123op dbname=dzikra_pos sslmode=disable" up

goose-down:
# example : make goose-down
	@echo " >> Installing goose if not installed"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@goose -dir db/migrations postgres "host=localhost user=postgres password=21012123op dbname=dzikra_pos sslmode=disable" down

goose-status:
# example : make goose-status
	@echo " >> Installing goose if not installed"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@goose -dir db/migrations postgres "host=localhost user=postgres password=21012123op dbname=dzikra_pos sslmode=disable" status

seed:
# make seed total=10 table=roles
	$(GO_CMD) run $(MAIN) seed -total=$(total) -table=$(table)

PROTO_SRC_DIR := ./cmd/proto/tokenvalidation
PROTO_OUT_DIR := ./cmd/proto/tokenvalidation
PROTO_FILE := tokenvalidation.proto

generate-proto:
	protoc --proto_path=$(PROTO_SRC_DIR) \
		--go_out=$(PROTO_OUT_DIR) --go-grpc_out=$(PROTO_OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_SRC_DIR)/$(PROTO_FILE)

docker-build:
	@echo " >> Building Docker image"
	docker buildx build --platform linux/amd64 -t ghcr.io/hilmiikhsan/dzikra-user-service:latest --push .

wire:
	@echo ">> Running Wire in internal/module/user/handler/rest"
	cd internal/module/user/handler/rest && wire
