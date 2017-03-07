DEFAULT_NAME := slackbotwit

package = github.com/lenfree/$(DEFAULT_NAME)

all: release

.PHONY: install
install:
	go get -v

.PHONY: script
script:
	go vet ./...
	go build -v -race ./...

.PHONY: release
release: install script
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/$(DEFAULT_NAME)-linux-amd64 $(package)
	GOOS=linux GOARCH=arm go build -o release/$(DEFAULT_NAME)-linux-arm $(package)
	GOOS=darwin GOARCH=amd64 go build -o release/$(DEFAULT_NAME)-darwin-amd64 $(package)

.PHONY: run
run:
	@go run **.go

.PHONY: buildpi
buildpi: checkdep
	GOOS=linux GOARCH=arm GOARM=6 go build -o release/$(DEFAULT_NAME)-arm-$(VERSION)

.PHONY: runpi
runpi:
	nohup ./$(DEFAULT_NAME) &

.PHONY: checkdep
checkdep:
ifndef VERSION
	$(error VERSION is not set)
endif