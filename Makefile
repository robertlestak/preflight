bin: bin/preflight_darwin_amd64 bin/preflight_linux_amd64 bin/preflight_windows_amd64.exe
bin: bin/preflight_darwin_arm64 bin/preflight_linux_arm64 bin/preflight_windows_arm64.exe

bin/preflight_darwin_amd64:
	@mkdir -p bin
	@echo "Compiling preflight..."
	GOOS=darwin GOARCH=amd64 go build -o $@ cmd/preflight/*.go

bin/preflight_darwin_arm64:
	@mkdir -p bin
	@echo "Compiling preflight..."
	GOOS=darwin GOARCH=arm64 go build -o $@ cmd/preflight/*.go

bin/preflight_linux_amd64:
	@mkdir -p bin
	@echo "Compiling preflight..."
	GOOS=linux GOARCH=amd64 go build -o $@ cmd/preflight/*.go

bin/preflight_linux_arm64:
	@mkdir -p bin
	@echo "Compiling preflight..."
	GOOS=linux GOARCH=arm64 go build -o $@ cmd/preflight/*.go

bin/preflight_windows_amd64.exe:
	@mkdir -p bin
	@echo "Compiling preflight..."
	GOOS=windows GOARCH=amd64 go build -o $@ cmd/preflight/*.go

bin/preflight_windows_arm64.exe:
	@mkdir -p bin
	@echo "Compiling preflight..."
	GOOS=windows GOARCH=arm64 go build -o $@ cmd/preflight/*.go

.PHONY: install
install: bin
	@echo "Installing preflight..."
	@scp bin/preflight_$$(go env GOOS)_$$(go env GOARCH) /usr/local/bin/preflight