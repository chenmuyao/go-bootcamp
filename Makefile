all: dev

run: dev
	@./webook

test:
	@mockgen -source=./internal/service/user.go -package=mock -destination=./internal/service/mocks/user.mock.go
	@mockgen -source=./internal/service/code.go -package=mock -destination=./internal/service/mocks/code.mock.go
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
