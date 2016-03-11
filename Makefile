.DEFAULT_GOAL := build

.PHONY: init
init:
	glide update  --no-recursive
	git init
	git add .

.PHONY: update-deps
update-deps:
	glide update  --no-recursive
    
.PHONY: test
test:
	go test $(glide novendor)

.PHONY: build
build:
	go build -v -o twitterbeat
