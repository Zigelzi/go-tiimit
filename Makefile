# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
dev/templ:
	@echo "\n\nStarting templ watch"
	templ generate -path ./cmd/web -watch -proxy="http://localhost:8080" --open-browser=false

# run air to detect any go file changes to re-build and re-run the server.
dev/server:
	@echo "\n\nStarting Air to watch server"
	air

# run tailwindcss to generate the styles.css bundle in watch mode.
dev/tailwind:
	@echo "\n\nStarting Tailwind watcher"
	npx --yes @tailwindcss/cli -i ./cmd/web/tailwind.css -o ./cmd/web/static/tailwind.css --minify --watch


# watch for any js or css change in the static/ folder, then reload the browser via templ proxy.
dev/sync_assets:
	@echo "\n\nStarting sync assets"
	air \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "./cmd/web/static" \
	--build.include_ext "css"

# start all 5 watch processes in parallel.
dev: 
	make -j4 dev/tailwind dev/server dev/templ dev/sync_assets

prod/build-arm64:
	make prod/tailwind
	make prod/build-server-arm64
	
prod/tailwind:
	npx tailwindcss -i ./cmd/web/tailwind.css -o ./cmd/web/static/tailwind.css --minify

prod/build-server-arm64:
	templ generate
	mkdir -p build
	GOOS=linux GOARCH=arm64 go build -o ./build ./cmd/web
