.PHONY: test

TEST_ARGS = -vet=all -failfast -fullpath -cover -race -timeout=$(if $(timeout), $(count),5s) $(if $(count),-count=$(count)) $(if $(run),-run=$(run)) $(if $(package),$(package),./...)

.PHONY: test
test:
	go test $(TEST_ARGS)

.PHONY: testsum
testsum:
	gotestsum -- $(TEST_ARGS)

mockgen\:install:
	go install github.com/golang/mock/mockgen@v1.6.0

generate:
	go generate ./...

deadcode\:install:
	go install golang.org/x/tools/cmd/deadcode@latest

clean\:testcache:
	go clean -testcache

run:
	go run ./cmd