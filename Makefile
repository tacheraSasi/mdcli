VERSION=1

build:
	make build_android
	make build_linux
	make build_mac
	make build_windows

build_linux:
	@echo 'building linux binary...'
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o mdcli
	@echo 'shrinking binary...'
	./upx --brute mdcli
	@echo 'zipping build...'
	tar -zcvf bin/mdcli_linux_amd64.tar.gz mdcli
	@echo 'cleaning up...'
	rm mdcli

build_windows:
	@echo 'building windows executable...'
	env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o mdcli_windows_amd64.exe
	@echo 'shrinking build...'
	./upx --brute bin/mdcli_windows_amd64.exe

build_mac:
	@echo 'building mac binary...'
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o mdcli
	@echo 'shrinking binary...'
	./upx --brute mdcli
	@echo 'zipping build...'
	tar -zcvf bin/mdcli_mac_amd64.tar.gz mdcli
	@echo 'cleaning up...'
	rm mdcli

build_android:
	@echo 'building android binary'
	env GOOS=android GOARCH=arm64 go build -ldflags="-s -w" -o mdcli
	@echo 'zipping build...'
	tar -zcvf bin/mdcli_android_arm64.tar.gz mdcli
	@echo 'cleaning up...'
	rm mdcli

build_test:
	go build -ldflags="-s -w" -o mdcli

dependencies:
	@echo 'checking dependencies...'
	go mod tidy

clean:
	go clean
