VERSION=$(shell git describe --tags)

build:
	go build -ldflags "-s -w -X github.com/lemonade-command/lemonade/lemon.Version=$(VERSION)"

install:
	go install -ldflags "-s -w -X github.com/lemonade-command/lemonade/lemon.Version=$(VERSION)"

release:
	CGO_ENABLED=0 gox --arch 'amd64 arm64' --os 'windows linux darwin' --output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}/{{.Dir}}" -ldflags "-s -w -X github.com/lemonade-command/lemonade/lemon.Version=$(VERSION)"
	zip      pkg/lemonade_windows_amd64.zip   dist/lemonade_windows_amd64/lemonade.exe -j
	tar zcvf pkg/lemonade_linux_amd64.tar.gz  -C dist/lemonade_linux_amd64/  lemonade
	tar zcvf pkg/lemonade_linux_arm64.tar.gz  -C dist/lemonade_linux_arm64/  lemonade
	tar zcvf pkg/lemonade_darwin_amd64.tar.gz -C dist/lemonade_darwin_amd64/ lemonade
	tar zcvf pkg/lemonade_darwin_arm64.tar.gz -C dist/lemonade_darwin_arm64/ lemonade

clean:
	rm -rf dist/
	rm -f pkg/*.tar.gz pkg/*.zip
