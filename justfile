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

# TODO: Add tests
