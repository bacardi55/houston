GOCMD := CGO_ENABLED=0 go
BINARY := houston
BINDIR := ./bin
VERSION := 1.0.0

GOLDFLAGS := -s -w -X main.Version=$(VERSION)

BUILD_TIME := ${shell date "+%Y-%m-%dT%H:%M"}

.PHONY: build
build:
	${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}

.PHONY: clean
clean:
	rm -f ${BINDIR}/${BINARY}

fmt:
	go fmt ./...

.PHONY: release
release:
	echo "Tagging version ${VERSION}"
	git tag -a v${VERSION} -m "New released tag: v${VERSION}"
	GOOS=linux GOARCH=amd64 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_arm
	GOOS=linux GOARCH=arm64 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_arm64
	GOOS=linux GOARCH=386 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_386

.PHONY: dependencies
dependencies:
	${GOCMD} get "git.sr.ht/~adnano/go-gemini"
