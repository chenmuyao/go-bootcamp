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
	@./wetravel

.PHONY: cover
cover:
	@gocov test ./... | gocov-html -t kit > ./tmp/report.html && firefox ./tmp/report.html

.PHONY: mock
mock:
	@mockgen -source=./internal/service/user.go -package=svcmocks -destination=./internal/service/mocks/user.mock.go
	@mockgen -source=./internal/service/code.go -package=svcmocks -destination=./internal/service/mocks/code.mock.go
	@mockgen -source=./internal/service/sms/types.go -package=smsmocks -destination=./internal/service/sms/mocks/sms.mock.go
	@mockgen -source=./internal/service/article.go -package=svcmocks -destination=./internal/service/mocks/article.mock.go
	@mockgen -source=./internal/repository/user.go -package=repomocks -destination=./internal/repository/mocks/user.mock.go
	@mockgen -source=./internal/repository/code.go -package=repomocks -destination=./internal/repository/mocks/code.mock.go
	@mockgen -source=./internal/repository/sms.go -package=repomocks -destination=./internal/repository/mocks/sms.mock.go
	@mockgen -source=./internal/repository/article.go -package=repomocks -destination=./internal/repository/mocks/article.mock.go
	@mockgen -source=./internal/repository/article_author.go -package=repomocks -destination=./internal/repository/mocks/article_author.mock.go
	@mockgen -source=./internal/repository/article_reader.go -package=repomocks -destination=./internal/repository/mocks/article_reader.mock.go
	@mockgen -source=./internal/repository/dao/user.go -package=daomocks -destination=./internal/repository/dao/mocks/user.mock.go
	@mockgen -source=./internal/repository/dao/sms.go -package=daomocks -destination=./internal/repository/dao/mocks/sms.mock.go
	@mockgen -source=./internal/repository/dao/article.go -package=daomocks -destination=./internal/repository/dao/mocks/article.mock.go
	@mockgen -source=./internal/repository/dao/article_author.go -package=daomocks -destination=./internal/repository/dao/mocks/article_author.mock.go
	@mockgen -source=./internal/repository/dao/article_reader.go -package=daomocks -destination=./internal/repository/dao/mocks/article_reader.mock.go
	@mockgen -source=./internal/repository/cache/types.go -package=cachemocks -destination=./internal/repository/cache/mocks/cache.mock.go
	@mockgen -source=./pkg/limiter/types.go -package=limitermocks -destination=./pkg/limiter/mocks/limiter.mock.go
	@mockgen -package=redismock -destination=./internal/repository/cache/rediscache/mocks/redismock.mock.go github.com/redis/go-redis/v9 Cmdable
	@mockgen -source=./internal/web/jwt/types.go -package=jwtmocks -destination=./internal/web/jwt/mocks/jwt.mock.go
	@cd ./internal/integration/startup/ && wire && cd -

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
