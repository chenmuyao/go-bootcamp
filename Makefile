all: dev

run: dev
	@./webook

.PHONY: dev
dev:
	@rm -f webook
	@go mod tidy
	@go build -o webook .

.PHONY: docker
docker:
	@rm -f webook
	@go mod tidy
	@GOOS=linux GOARCH=arm go build --tags=k8s -o webook .
	@docker rmi -f vinchent123/webook:v0.0.1
	@docker build -t vinchent123/webook:v0.0.1 .
