export GO111MODULE=on

.PHONY: all
all: test vet lint fmt

.PHONY: test
test:
	@go test ./gettr -cover

.PHONY: vet
vet:
	@go vet -all ./gettr

.PHONY: lint
lint:
	@golint -set_exit_status ./...

.PHONY: fmt
fmt:
	@test -z $$(go fmt ./...)

