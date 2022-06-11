VERSION = $(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"
OSARCH=$(shell go env GOHOSTOS)-$(shell go env GOHOSTARCH)

NANOSTARTER=\
	nanostarter-linux-amd64 \
	nanostarter-darwin-amd64 \
	nanostarter-darwin-arm64

my: nanostarter-$(OSARCH)

$(NANOSTARTER):
	GOOS=$(word 2,$(subst -, ,$@)) GOARCH=$(word 3,$(subst -, ,$(subst .exe,,$@))) go build $(LDFLAGS) -o $@ ./$<

clean:
	rm -f nanostarter-*

test:
	go test -v -cover -race ./...

