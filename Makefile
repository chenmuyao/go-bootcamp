# Dependencies :
# wire
#
# docker
# mockgen
#
# buf -- to handle protobuf
#
# NOTE: firefox can only be run on local machine. Need other commands to run
# on a CI pipeline
# gocov
# gocov-html
# firefox

all: dev

.PHONY: run
run: dev
	@./wetravel

.PHONY: cover
cover:
	@gocov test ./... | gocov-html -t kit > ./tmp/report.html && firefox ./tmp/report.html

.PHONY: mock
mock:
	@go generate ./...
	@mockgen -package=redismock -destination=./internal/repository/cache/rediscache/mocks/redismock.mock.go github.com/redis/go-redis/v9 Cmdable
	@mockgen -package=intrv1mock -source=./api/proto/gen/intr/v1/interactive_grpc.pb.go -destination=./api/proto/gen/intr/v1/mock/intrv1mock.mock.go
	@cd ./internal/integration/startup/ && wire && cd -
	@cd ./interactive/integration/startup/ && wire && cd -

.PHONY: test
test: mock
	@go test ./...

.PHONY: dev
dev:
	@rm -f wetravel
	@go mod tidy
	@wire
	@go build -v -o wetravel .

.PHONY: docker
docker:
	@rm -f wetravel
	@go mod tidy
	# @GOOS=linux GOARCH=arm go build --tags=k8s -o wetravel .
	@GOOS=linux GOARCH=arm go build -o wetravel .
	@docker rmi -f vinchent123/wetravel:v0.0.1
	@docker build -t vinchent123/wetravel:v0.0.1 .

.PHONY: grpc
grpc:
	@npx buf generate api/proto/
