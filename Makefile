GO111MODULES=on
APP?=prometheus-tester
REGISTRY?=masonmei
COMMIT_SHA=$(shell git rev-parse --short HEAD)
GOOS:=linux
GOARCH:=amd64


prometheus-bm: cmd/prometheus-bm/prometheus-bm

cmd/prometheus-bm/prometheus-bm: cmd/prometheus-bm/main.go
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o $@ ./$(@D)

.PHONY: build
## build: build the application
build: clean prometheus-bm
	@echo "Building..."

.PHONY: run
## run: runs go run main.go
run:
	go run -race main.go

.PHONY: clean
## clean: cleans the binary
clean:
	@echo "Cleaning"
	@rm -rf cmd/prometheus-bm/prometheus-bm
	@go clean ./...

.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...


.PHONY: build-tokenizer
## build-tokenizer: build the tokenizer application
build-tokenizer:
	${MAKE} -c tokenizer build

.PHONY: setup
## setup: setup go modules
setup:
	@go mod init \
		&& go mod tidy \
		&& go mod vendor

# helper rule for deployment
#check-environment:
#ifndef ENV
#    $(error ENV not set, allowed values - `staging` or `production`)
#endif

.PHONY: docker-build
## docker-build: builds the stringifier docker image to registry
docker-build: build
	docker build -t ${REGISTRY}/${APP}:${IMAGE_TAG} -f Dockerfile .

.PHONY: docker-push
## docker-push: pushes the stringifier docker image to registry
docker-push: prometheus-bm docker-build
	docker push ${REGISTRY}/${APP}:${COMMIT_SHA}

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
