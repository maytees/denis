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
