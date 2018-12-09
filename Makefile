VERSION=$(shell git describe --tags)

build:
	go build -ldflags "-X github.com/lemonade-command/lemonade/lemon.Version=$(VERSION)"

install:
	go install -ldflags "-X github.com/lemonade-command/lemonade/lemon.Version=$(VERSION)"

release:
	gox --arch 'amd64 386' --os 'windows linux darwin' --output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}/{{.Dir}}" -ldflags "-X github.com/lemonade-command/lemonade/lemon.Version=$(VERSION)"
	zip      pkg/lemonade_windows_386.zip     dist/lemonade_windows_386/lemonade.exe   -j
	zip      pkg/lemonade_windows_amd64.zip   dist/lemonade_windows_amd64/lemonade.exe -j
	tar zcvf pkg/lemonade_linux_386.tar.gz    -C dist/lemonade_linux_386/    lemonade
	tar zcvf pkg/lemonade_linux_amd64.tar.gz  -C dist/lemonade_linux_amd64/  lemonade
	tar zcvf pkg/lemonade_darwin_386.tar.gz   -C dist/lemonade_darwin_386/   lemonade
	tar zcvf pkg/lemonade_darwin_amd64.tar.gz -C dist/lemonade_darwin_amd64/ lemonade

clean:
	rm -rf dist/
	rm -f pkg/*.tar.gz pkg/*.zip
