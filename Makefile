VERSION = $(shell git describe --tags --always --dirty)
#go tool nm ./nanostarter-linux-amd64 | grep version #to get the path of the version variable
#https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
LDFLAGS=-ldflags "-X 'github.com/sheshenia/nanostarter/server.version=$(VERSION)'"
OSARCH=$(shell go env GOHOSTOS)-$(shell go env GOHOSTARCH)

NANOSTARTER=\
	nanostarter-linux-amd64 \
	nanostarter-darwin-amd64 \
	nanostarter-darwin-arm64

my: nanostarter-$(OSARCH)

$(NANOSTARTER):
	cd ./client && npm run build
	rm -rf ./server/client
	mv ./client/dist ./server/client #move to server folder for golang embedding
	GOOS=$(word 2,$(subst -, ,$@)) GOARCH=$(word 3,$(subst -, ,$(subst .exe,,$@))) go build $(LDFLAGS) -o $@ ./$<

# dev mode client
dev_c:
	cd ./client && npm run dev

# dev mode server
dev_s:
	go run . -embed false

clean:
	rm -rf nanostarter-*

test:
	go test -v -cover -race ./...

nanostarter-%-$(VERSION).zip: nanostarter-%
	# creating folder with GOOS-VERSION
	mkdir -p nanostarter-$(*)-$(VERSION)
	# move binary to created folder $(*) = % in command
	mv nanostarter-$(*) nanostarter-$(*)-$(VERSION)
	# copy release template
	cp -r ./internal/release_template/* ./nanostarter-$(*)-$(VERSION)
	# replace command examples to correct GOOS and GOARCH
	sed -i -e 's/linux-amd64/$(*)/g' ./nanostarter-$(*)-$(VERSION)/BEFORE_START.md
	# move release folder structure to archive
	cd ./nanostarter-$(*)-$(VERSION) && zip -m ../$@ *
	# delete release folder
	rm -r ./nanostarter-$(*)-$(VERSION)

release: \
	nanostarter-linux-amd64-$(VERSION).zip \
	nanostarter-darwin-amd64-$(VERSION).zip \
	nanostarter-darwin-arm64-$(VERSION).zip

.PHONY: my $(NANOSTARTER) clean release test dev_c dev_s