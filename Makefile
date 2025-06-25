include .envrc
MIGRATIONS_PATH = ./cmd/migrate/migrations


.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-force
migrate-force:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) force $(filter-out $@,$(MAKECMDGOALS))


templ:
	templ generate --watch --proxy="http://localhost:8090" --open-browser=false

server:
	air \
	--build.cmd "go build -o tmp/bin/main ./cmd/myportfolio/main.go" \
	--build.bin "tmp/bin/main" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

tailwind:
	npx @tailwindcss/cli -i ./assets/css/input.css -o ./assets/css/output.css --watch
	

dev:
	make -j3 tailwind templ server


