
build-docker:
	docker build --no-cache --rm --tag tmaize/webdav:1.0.0 .
	docker image prune -f
clean:
	rm -rf dist
	mkdir -p dist
build: clean build-linux-amd64 build-windows-amd64 build-mac-amd64 build-mac-arm64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o dist/webdav .
	cd dist && 7z a -sdel webdav-linux-amd64.tar webdav
	cd dist && 7z a -sdel webdav-linux-amd64.tar.gz webdav-linux-amd64.tar
build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o dist/webdav.exe .
	cd dist && 7z a -sdel webdav-windows-amd64.zip webdav.exe
build-mac-amd64:
	GOOS=darwin GOARCH=amd64 go build -o dist/webdav .
	cd dist && 7z a -sdel webdav-darwin-amd64.tar webdav
	cd dist && 7z a -sdel webdav-darwin-amd64.tar.gz webdav-darwin-amd64.tar
build-mac-arm64:
	GOOS=darwin GOARCH=arm64 go build -o dist/webdav .
	cd dist && 7z a -sdel webdav-darwin-arm64.tar webdav
	cd dist && 7z a -sdel webdav-darwin-arm64.tar.gz webdav-darwin-arm64.tar
