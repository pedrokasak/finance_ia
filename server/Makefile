.PHONY: run build test test-verbose coverage coverage-html lint docker-up docker-down docker-logs migrate clean dev

COVERAGE_FILE=coverage.out
COVERAGE_THRESHOLD=90

# ── Dev ───────────────────────────────────────────────────────
run:
	go run main.go

dev:
	air # hot reload — install: go install github.com/cosmtrek/air@latest

# ── Tests ─────────────────────────────────────────────────────
test:
	go test ./... -race -count=1 -timeout 60s

test-verbose:
	go test ./... -race -v -count=1 -timeout 60s

coverage:
	go test ./... -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic -timeout 60s
	@TOTAL=$$(go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Coverage: $${TOTAL}%"; \
	if [ $$(echo "$${TOTAL} < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Below $(COVERAGE_THRESHOLD)% threshold"; exit 1; \
	else \
		echo "✅ Meets $(COVERAGE_THRESHOLD)% threshold"; \
	fi

coverage-html: coverage
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	open coverage.html

# ── Build ──────────────────────────────────────────────────────
build:
	go build -ldflags="-s -w" -o bin/server ./main.go

# ── Lint ───────────────────────────────────────────────────────
lint:
	go vet ./...
	golangci-lint run ./... --timeout 5m

install-deps:
	go mod tidy && go mod download
	@which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@which air || go install github.com/cosmtrek/air@latest

# ── Docker ─────────────────────────────────────────────────────
docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

# ── Newman (Postman) ───────────────────────────────────────────
postman-test:
	npx newman run postman/finance-ia.postman_collection.json \
		--environment postman/local.postman_environment.json \
		--reporters cli,html --reporter-html-export postman-report.html

# ── DB ─────────────────────────────────────────────────────────
migrate:
	go run main.go -migrate-only

# ── Cleanup ────────────────────────────────────────────────────
clean:
	rm -rf bin/ $(COVERAGE_FILE) coverage.html postman-report.html
