.PHONY: new daily test cover lint fmt vet tidy

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
	gofmt -l -w .

vet:
	go vet ./...

lint: fmt vet

tidy:
	go mod tidy
