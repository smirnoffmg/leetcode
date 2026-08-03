.PHONY: new daily submit test cover bench lint fmt vet tidy hooks push

# Куки leetcode, если лежат в .env (файл не под гитом)
-include .env
export LEETCODE_SESSION LEETCODE_CSRF

# make new S=two-sum  (N/T/D — фолбэк, если leetcode недоступен)
new:
	@./scripts/new.sh $(S) $(N) "$(T)" $(D)

daily:
	@./scripts/daily.sh

# make submit            → отправить задачу, которую правил последней
# make submit S=two-sum  → отправить конкретную
# нужны LEETCODE_SESSION и LEETCODE_CSRF из cookies браузера
submit:
	@go run ./cmd/lcsubmit -slug "$(S)"

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

# make push            → solve(2026-07-25): 3536-maximum-product-of-two-digits
# make push M="..."    → свой текст коммита
push:
	@git add -A
	@if git diff --cached --quiet; then \
		echo "нечего коммитить"; exit 0; \
	fi; \
	msg="$(M)"; \
	if [ -z "$$msg" ]; then \
		day=$$(date +%F); \
		slugs=$$(git diff --cached --name-only -- problems \
			| cut -d/ -f2 | sort -u | paste -sd, - | sed 's/,/, /g'); \
		if [ -n "$$slugs" ]; then msg="solve($$day): $$slugs"; \
		else msg="chore($$day): update"; fi; \
	fi; \
	git commit -m "$$msg" && git pull --rebase && git push
