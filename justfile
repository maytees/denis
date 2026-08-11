# Run with the air hot reloading
[group('dev')]
dev:
	sudo air

# Run without hot reloading
[group('dev')]
dev-run:
	sudo go run .

[group('dev')]
build:
	go build -o dist/denis

# Run without hot reloading
[group('dev')]
start: build
	sudo dist/denis

# Run all tests
[group('dev')]
test:
	go test ./...

# Use gotestsum to test
[group('dev')]
prettytest:
	gotestsum --format testdox

# Clean build files
[group('dev')]
clean:
	rm -rf dist

# Put example config/records in ./config
[group('dev')]
config-example:
	cp ./config/config.example.toml ./config/config.toml
	cp ./config/records.example.toml ./config/records.toml

# Lint with golangci-lint
[group('dev')]
lint:
	golangci-lint run

# Format with golangci-lint
[group('dev')]
format:
	golangci-lint fmt

# Build the image and (re)start the container
[group('docker')]
docker-compose:
	docker compose up -d --build

# Follow container logs
[group('docker')]
docker-logs:
	docker compose logs -f

# Restart the container (e.g. after editing records.toml)
[group('docker')]
docker-restart:
	docker compose restart

# Stop and remove the container
[group('docker')]
docker-down:
	docker compose down

# Verify DENIS answers: local records + upstream forwarding
[group('docker')]
docker-check:
	dig @127.0.0.1 localhost
	dig @127.0.0.1 google.com

# Point macOS DNS at DENIS
[group('dns')]
dns-use:
	sudo networksetup -setdnsservers "Wi-Fi" 127.0.0.1

# Revert macOS DNS to DHCP (escape hatch if DENIS is down)
[group('dns')]
dns-revert:
	sudo networksetup -setdnsservers "Wi-Fi" "Empty"

# Run the docs site tailwind watcher
[group('site')]
site-css:
	cd www && bun run css

# Run the docs site dev server
[group('site')]
site-dev:
	cd www && bun run dev

# Build the docs site to www/dist
[group('site')]
site-build:
	cd www && bun install && bun run build

# Update the vendored bloom runtime from a sibling ../bloom checkout
[group('site')]
site-vendor:
	rsync -a --delete --exclude='*.test.ts' --exclude='tsconfig.json' --exclude='README.md' ../bloom/packages/runtime/ www/vendor/bloom/
	cd www && bun install

# TODO: Add tests
