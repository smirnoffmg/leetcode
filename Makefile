.PHONY: new daily test cover bench lint fmt vet tidy hooks

# make new S=two-sum  (N/T/D — фолбэк, если leetcode недоступен)
new:
	@./scripts/new.sh $(S) $(N) "$(T)" $(D)

daily:
	@./scripts/daily.sh

test:
	go test ./...

cover:
	go test -cover ./...

bench:
	go test -bench=. -benchmem ./...

fmt:
	golangci-lint fmt

vet:
	go vet ./...

lint:
	golangci-lint run

hooks:
	pre-commit install
	pre-commit run --all-files

tidy:
	go mod tidy
