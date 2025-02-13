VERSION=0.0.1

build:
	make build_android
	make build_linux
	make build_mac
	make build_windows

build_linux:
	@echo 'building linux binary...'
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ellie
	@echo 'shrinking binary...'
	./upx --brute ellie
	@echo 'zipping build...'
	tar -zcvf binaries/ellie_linux_amd64.tar.gz ellie
	@echo 'cleaning up...'
	rm ellie

build_windows:
	@echo 'building windows executable...'
	env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ellie_windows_amd64.exe
	@echo 'shrinking build...'
	./upx --brute binaries/ellie_windows_amd64.exe

build_mac:
	@echo 'building mac binary...'
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ellie
	@echo 'shrinking binary...'
	./upx --brute ellie
	@echo 'zipping build...'
	tar -zcvf binaries/ellie_mac_amd64.tar.gz ellie
	@echo 'cleaning up...'
	rm ellie

build_android:
	@echo 'building android binary'
	env GOOS=android GOARCH=arm64 go build -ldflags="-s -w" -o ellie
	@echo 'zipping build...'
	tar -zcvf binaries/ellie_android_arm64.tar.gz ellie
	@echo 'cleaning up...'
	rm ellie

build_test:
	go build -ldflags="-s -w" -o ellie

dependencies:
	@echo 'checking dependencies...'
	go mod tidy



clean:
	go clean
