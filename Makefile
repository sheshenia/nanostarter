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
	rm -f nanostarter-*

test:
	go test -v -cover -race ./...

