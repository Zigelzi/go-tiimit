# Tiimit build commands.
#
# `make generate` runs all code generation: sqlc (SQL -> Go), templ
# (.templ -> Go) and the standalone Tailwind CLI (input CSS -> the stylesheet
# the web binary embeds). `run` and the builds depend on it, so a fresh
# checkout works with one command. `test` needs only the templ code.

TAILWIND_IN  = ./cmd/web/tailwind.css
TAILWIND_OUT = ./cmd/web/static/tailwind.css

# The standalone Tailwind CLI is downloaded once into ./bin (gitignored) and
# checked against the release's published SHA-256. The version is the one the
# bun build used, so the stylesheet stays byte-for-byte the same. To upgrade,
# change the version and copy the four sums from the release's sha256sums.txt.
TAILWIND_VERSION = 4.1.17
TAILWIND_SHA256_linux-x64   = cc115d9b6c4ede4e423bfea6a3cfc3f03e6c1702b7d910369b9540e2b4cf3860
TAILWIND_SHA256_linux-arm64 = 264689984311cc2ca104a9d6a4d039b4267aecf4543dc9ea0bbc094eeb0ccc8d
TAILWIND_SHA256_macos-x64   = 6337c87185321c079244dfbd9c24502a3006401bd1539d19cdc9df8de9101843
TAILWIND_SHA256_macos-arm64 = 852b62d91a6dfa56969686bc8a4bae0e4722de8f5a9b42b102bb589b7ece8cfe

TAILWIND_OS       := $(if $(filter Darwin,$(shell uname -s)),macos,linux)
TAILWIND_ARCH     := $(or $(if $(filter x86_64 amd64,$(shell uname -m)),x64),$(if $(filter arm64 aarch64,$(shell uname -m)),arm64),$(shell uname -m))
TAILWIND_PLATFORM := $(TAILWIND_OS)-$(TAILWIND_ARCH)
TAILWIND          := ./bin/tailwindcss-$(TAILWIND_VERSION)

.PHONY: generate run test dev dev/templ dev/server dev/tailwind prod/build-arm64 deploy

generate: $(TAILWIND)
	sqlc generate
	templ generate
	$(TAILWIND) -i $(TAILWIND_IN) -o $(TAILWIND_OUT) --minify

# run starts the server with the settings in .env.
run: generate
	go run ./cmd/web

# test generates the templ code the tests compile against, then runs every test.
test:
	templ generate
	go test ./...

# prod/build-arm64 builds the binary for the Raspberry Pi into ./build.
prod/build-arm64: generate
	mkdir -p build
	GOOS=linux GOARCH=arm64 go build -o ./build ./cmd/web

# deploy tests, builds and installs Tiumo on the Pi, backing up the database
# first. A deploy is also a migration of the live database. Needs
# deploy/deploy.conf; docs/deployment-guide.md has the setup.
deploy:
	./deploy/deploy.sh

# dev runs three watchers in parallel. Open http://localhost:7331.
#   templ:    regenerates templates and proxies to the server with live reload
#   air:      rebuilds and restarts the server when Go code or CSS changes
#   tailwind: rebuilds the stylesheet when classes change
# After editing SQL, run `make generate` in another terminal: sqlc isn't watched.
dev: $(TAILWIND)
	$(MAKE) -j3 dev/tailwind dev/server dev/templ

dev/templ:
	@echo "\n\nStarting templ watch"
	templ generate -path ./cmd/web -watch -proxy="http://localhost:8080" --open-browser=false

dev/server:
	@echo "\n\nStarting Air to watch server"
	air

dev/tailwind:
	@echo "\n\nStarting Tailwind watcher"
	$(TAILWIND) -i $(TAILWIND_IN) -o $(TAILWIND_OUT) --minify --watch

$(TAILWIND):
	@if [ -z "$(TAILWIND_SHA256_$(TAILWIND_PLATFORM))" ]; then \
		echo "No pinned Tailwind CLI for $(TAILWIND_PLATFORM): add it to the Makefile from the release's sha256sums.txt"; \
		exit 1; \
	fi
	mkdir -p bin
	curl -sSLf -o $@.download \
		https://github.com/tailwindlabs/tailwindcss/releases/download/v$(TAILWIND_VERSION)/tailwindcss-$(TAILWIND_PLATFORM)
	@sum=$$( (sha256sum $@.download 2>/dev/null || shasum -a 256 $@.download) | cut -d' ' -f1 ); \
	if [ "$$sum" != "$(TAILWIND_SHA256_$(TAILWIND_PLATFORM))" ]; then \
		echo "Tailwind download's SHA-256 is $$sum, want $(TAILWIND_SHA256_$(TAILWIND_PLATFORM))"; \
		rm -f $@.download; exit 1; \
	fi
	chmod +x $@.download
	mv $@.download $@
