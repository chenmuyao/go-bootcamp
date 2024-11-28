all: dev

.PHONY: run
run: dev
	@./webook

.PHONY: cover
cover:
	@go test ./... -coverprofile ./tmp/cover.out && go tool cover -html ./tmp/cover.out -o ./tmp/cover.html && firefox ./tmp/cover.html

.PHONY: test
test:
	@mockgen -source=./internal/service/user.go -package=svcmocks -destination=./internal/service/mocks/user.mock.go
	@mockgen -source=./internal/service/code.go -package=svcmocks -destination=./internal/service/mocks/code.mock.go
	@mockgen -source=./internal/repository/user.go -package=repomocks -destination=./internal/repository/mocks/user.mock.go
	@mockgen -source=./internal/repository/code.go -package=repomocks -destination=./internal/repository/mocks/code.mock.go
	@mockgen -source=./internal/repository/dao/user.go -package=daomocks -destination=./internal/repository/dao/mocks/user.mock.go
	@mockgen -source=./internal/repository/cache/types.go -package=cachemocks -destination=./internal/repository/cache/mocks/cache.mock.go
	@mockgen -package=redismock -destination=./internal/repository/cache/rediscache/mocks/redismock.mock.go github.com/redis/go-redis/v9 Cmdable
	@cd ./internal/integration/startup/ && wire && cd -
	@go test -v ./...

.PHONY: dev
dev:
	@rm -f webook
	@go mod tidy
	@wire
	@go build -v -o webook .

.PHONY: docker
docker:
	@rm -f webook
	@go mod tidy
	@GOOS=linux GOARCH=arm go build --tags=k8s -o webook .
	@docker rmi -f vinchent123/webook:v0.0.1
	@docker build -t vinchent123/webook:v0.0.1 .
