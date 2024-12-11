# Dependencies :
# wire
#
# docker
# mockgen
#
# NOTE: firefox can only be run on local machine. Need other commands to run
# on a CI pipeline
# gocov
# gocov-html
# firefox

all: dev

.PHONY: run
run: dev
	@./webook

.PHONY: cover
cover:
	@gocov test ./... | gocov-html -t kit > ./tmp/report.html && firefox ./tmp/report.html

.PHONY: mock
mock:
	@mockgen -source=./internal/service/user.go -package=svcmocks -destination=./internal/service/mocks/user.mock.go
	@mockgen -source=./internal/service/code.go -package=svcmocks -destination=./internal/service/mocks/code.mock.go
	@mockgen -source=./internal/service/sms/types.go -package=smsmocks -destination=./internal/service/sms/mocks/sms.mock.go
	@mockgen -source=./internal/repository/user.go -package=repomocks -destination=./internal/repository/mocks/user.mock.go
	@mockgen -source=./internal/repository/code.go -package=repomocks -destination=./internal/repository/mocks/code.mock.go
	@mockgen -source=./internal/repository/sms.go -package=repomocks -destination=./internal/repository/mocks/sms.mock.go
	@mockgen -source=./internal/repository/dao/user.go -package=daomocks -destination=./internal/repository/dao/mocks/user.mock.go
	@mockgen -source=./internal/repository/dao/sms.go -package=daomocks -destination=./internal/repository/dao/mocks/sms.mock.go
	@mockgen -source=./internal/repository/cache/types.go -package=cachemocks -destination=./internal/repository/cache/mocks/cache.mock.go
	@mockgen -source=./pkg/limiter/types.go -package=limitermocks -destination=./pkg/limiter/mocks/limiter.mock.go
	@mockgen -package=redismock -destination=./internal/repository/cache/rediscache/mocks/redismock.mock.go github.com/redis/go-redis/v9 Cmdable
	@cd ./internal/integration/startup/ && wire && cd -

.PHONY: test
test: mock
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
