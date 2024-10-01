VERSION ?= 1

build:
	docker build \
	--progress plain \
	--no-cache \
	--tag codenotary-be:${VERSION} .

generate-mock:
	mockgen -source=./store/store.go -destination=./store/immudbvaultmock/mock.go -package=immudbvaultmock

test:
	go test ./... -count=1 -cover